// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package storage

import (
	"context"
	"fmt"

	"github.com/kk/kkartifact-server/internal/database"
)

// CleanupManager handles version cleanup
type CleanupManager struct {
	artifactManager *ArtifactManager
	db              *database.DB
}

// NewCleanupManager creates a new cleanup manager
func NewCleanupManager(artifactManager *ArtifactManager, db *database.DB) *CleanupManager {
	return &CleanupManager{
		artifactManager: artifactManager,
		db:              db,
	}
}

// CleanupOldVersions cleans up old versions beyond the retention limit
// This function deletes versions from both storage and database
func (cm *CleanupManager) CleanupOldVersions(ctx context.Context, project, app string, retentionLimit int) error {
	// Get project and app from database
	projectRepo := database.NewProjectRepository(cm.db)
	projectModel, err := projectRepo.CreateOrGet(project)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	appRepo := database.NewAppRepository(cm.db)
	appModel, err := appRepo.CreateOrGet(projectModel.ID, app)
	if err != nil {
		return fmt.Errorf("failed to get app: %w", err)
	}

	versionRepo := database.NewVersionRepository(cm.db)

	// Step 1: Clean up incomplete uploads (versions in storage without meta.yaml or without DB record)
	// This should be done first, before checking retention limits
	if err := cm.cleanupIncompleteVersions(ctx, project, app, appModel.ID, versionRepo); err != nil {
		// Log error but continue - this is not critical
		fmt.Printf("Warning: failed to cleanup incomplete versions for %s/%s: %v\n", project, app, err)
	}

	// Step 2: Clean up old versions beyond retention limit
	// Get all versions from database to count
	allVersions, err := versionRepo.ListByApp(appModel.ID, 10000, 0) // Get all versions
	if err != nil {
		return fmt.Errorf("failed to list versions from database: %w", err)
	}

	if len(allVersions) <= retentionLimit {
		return nil
	}

	// Calculate how many to delete
	toDelete := len(allVersions) - retentionLimit

	// Get oldest versions (from database, sorted by created_at ASC)
	oldestVersions, err := versionRepo.GetOldestVersions(appModel.ID, toDelete)
	if err != nil {
		return fmt.Errorf("failed to get oldest versions: %w", err)
	}

	// Delete from both storage and database
	for _, version := range oldestVersions {
		// Delete from storage first
		if err := cm.artifactManager.DeleteVersion(ctx, project, app, version.Hash); err != nil {
			// Log error but continue - storage might not exist
			fmt.Printf("Failed to delete version %s from storage: %v\n", version.Hash, err)
		}

		// Delete from database
		if err := versionRepo.Delete(appModel.ID, version.Hash); err != nil {
			fmt.Printf("Failed to delete version %s from database: %v\n", version.Hash, err)
			// Continue even if database delete fails
		}
	}

	return nil
}

// cleanupIncompleteVersions cleans up versions in storage that don't have meta.yaml or don't exist in database
func (cm *CleanupManager) cleanupIncompleteVersions(ctx context.Context, project, app string, appID int, versionRepo *database.VersionRepository) error {
	// Get all versions from storage
	storageVersions, err := cm.artifactManager.ListVersions(ctx, project, app)
	if err != nil {
		// If listing fails (e.g., app doesn't exist in storage), that's okay
		return nil
	}

	// Get all versions from database
	dbVersions, err := versionRepo.ListByApp(appID, 10000, 0)
	if err != nil {
		return fmt.Errorf("failed to list versions from database: %w", err)
	}

	// Create a map of database versions for quick lookup
	dbVersionMap := make(map[string]bool)
	for _, v := range dbVersions {
		dbVersionMap[v.Hash] = true
	}

	// Check each version in storage
	for _, version := range storageVersions {
		// Check if meta.yaml exists (indicates completed upload)
		// Use GetManifest to check if meta.yaml exists and is valid
		_, err := cm.artifactManager.GetManifest(ctx, project, app, version)
		hasManifest := err == nil // If GetManifest succeeds, meta.yaml exists

		// If meta.yaml doesn't exist (GetManifest failed), this is an incomplete upload - delete it
		if !hasManifest {
			fmt.Printf("Cleaning up incomplete version (no meta.yaml): %s/%s/%s\n", project, app, version)
			if err := cm.artifactManager.DeleteVersion(ctx, project, app, version); err != nil {
				fmt.Printf("Failed to delete incomplete version %s/%s/%s: %v\n", project, app, version, err)
			}
			continue
		}

		// If meta.yaml exists but version is not in database, delete it (orphaned version)
		if !dbVersionMap[version] {
			fmt.Printf("Cleaning up orphaned version (not in database): %s/%s/%s\n", project, app, version)
			if err := cm.artifactManager.DeleteVersion(ctx, project, app, version); err != nil {
				fmt.Printf("Failed to delete orphaned version %s/%s/%s: %v\n", project, app, version, err)
			}
		}
	}

	return nil
}
