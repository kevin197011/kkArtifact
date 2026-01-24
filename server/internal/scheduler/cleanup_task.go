// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package scheduler

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/kk/kkartifact-server/internal/database"
	"github.com/kk/kkartifact-server/internal/storage"
)

// CleanupTask is a scheduled task for version cleanup
type CleanupTask struct {
	db              *database.DB
	artifactManager *storage.ArtifactManager
	cleanupManager  *storage.CleanupManager
}

// NewCleanupTask creates a new cleanup task
func NewCleanupTask(db *database.DB, artifactManager *storage.ArtifactManager) *CleanupTask {
	cleanupManager := storage.NewCleanupManager(artifactManager, db)
	return &CleanupTask{
		db:              db,
		artifactManager: artifactManager,
		cleanupManager:  cleanupManager,
	}
}

// Name returns the task name
func (t *CleanupTask) Name() string {
	return "version-cleanup"
}

// Run runs the cleanup task
func (t *CleanupTask) Run(ctx context.Context) error {
	// Get retention limit from config
	configRepo := database.NewConfigRepository(t.db)
	retentionLimitStr, err := configRepo.Get("version_retention_limit")
	if err != nil {
		return fmt.Errorf("failed to get retention limit: %w", err)
	}

	retentionLimit, err := strconv.Atoi(retentionLimitStr)
	if err != nil {
		return fmt.Errorf("invalid retention limit: %w", err)
	}

	// Get all projects
	projectRepo := database.NewProjectRepository(t.db)
	projects, err := projectRepo.List(1000, 0) // Get up to 1000 projects
	if err != nil {
		return fmt.Errorf("failed to list projects: %w", err)
	}

	appRepo := database.NewAppRepository(t.db)

	// Iterate through all apps and cleanup
	for _, project := range projects {
		apps, err := appRepo.ListByProject(project.ID, 1000, 0)
		if err != nil {
			continue
		}

		for _, app := range apps {
			if err := t.cleanupManager.CleanupOldVersions(ctx, project.Name, app.Name, retentionLimit); err != nil {
				// Log error but continue
				log.Printf("Failed to cleanup versions for %s/%s: %v", project.Name, app.Name, err)
			}
		}
	}

	return nil
}

