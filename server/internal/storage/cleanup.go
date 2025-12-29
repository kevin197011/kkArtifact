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

	// Get all versions from database to count
	versionRepo := database.NewVersionRepository(cm.db)
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
