// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package services

import (
	"fmt"

	"github.com/kk/kkartifact-server/internal/database"
)

// Inventory represents the complete inventory structure
type Inventory struct {
	Projects []ProjectInventory `json:"projects"`
}

// ProjectInventory represents inventory data for a single project
type ProjectInventory struct {
	Project database.Project `json:"project"`
	Apps    []AppInventory   `json:"apps"`
}

// AppInventory represents inventory data for a single app
type AppInventory struct {
	App      database.App     `json:"app"`
	Versions []database.Version `json:"versions"`
}

// InventorySummary represents summary statistics for the inventory
type InventorySummary struct {
	TotalProjects int `json:"total_projects"`
	TotalApps     int `json:"total_apps"`
	TotalVersions int `json:"total_versions"`
}

// InventoryService provides methods to retrieve inventory data
type InventoryService struct {
	projectRepo *database.ProjectRepository
	appRepo     *database.AppRepository
	versionRepo *database.VersionRepository
}

// NewInventoryService creates a new inventory service
func NewInventoryService(
	projectRepo *database.ProjectRepository,
	appRepo *database.AppRepository,
	versionRepo *database.VersionRepository,
) *InventoryService {
	return &InventoryService{
		projectRepo: projectRepo,
		appRepo:     appRepo,
		versionRepo: versionRepo,
	}
}

// GetCompleteInventory retrieves the complete inventory structure with all projects, apps, and versions
func (s *InventoryService) GetCompleteInventory() (*Inventory, error) {
	// Fetch all projects with a large limit to get all projects
	projects, err := s.projectRepo.List(10000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects: %w", err)
	}

	inventory := &Inventory{
		Projects: make([]ProjectInventory, 0, len(projects)),
	}

	// For each project, fetch its apps and versions
	for _, project := range projects {
		projectInventory, err := s.buildProjectInventory(project)
		if err != nil {
			// Log error but continue with other projects
			// Return empty apps/versions for this project
			inventory.Projects = append(inventory.Projects, ProjectInventory{
				Project: *project,
				Apps:    []AppInventory{},
			})
			continue
		}
		inventory.Projects = append(inventory.Projects, *projectInventory)
	}

	return inventory, nil
}

// GetProjectInventory retrieves inventory data for a specific project
func (s *InventoryService) GetProjectInventory(projectName string) (*ProjectInventory, error) {
	// Get all projects and find by name
	projects, err := s.projectRepo.List(10000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects: %w", err)
	}

	var project *database.Project
	for _, p := range projects {
		if p.Name == projectName {
			project = p
			break
		}
	}

	if project == nil {
		return nil, fmt.Errorf("project not found: %s", projectName)
	}

	return s.buildProjectInventory(project)
}

// GetInventorySummary retrieves summary statistics for the inventory
func (s *InventoryService) GetInventorySummary() (*InventorySummary, error) {
	// Fetch all projects
	projects, err := s.projectRepo.List(10000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects: %w", err)
	}

	summary := &InventorySummary{
		TotalProjects: len(projects),
		TotalApps:     0,
		TotalVersions: 0,
	}

	// Count apps and versions for each project
	for _, project := range projects {
		apps, err := s.appRepo.ListByProject(project.ID, 10000, 0)
		if err != nil {
			// Continue with other projects even if one fails
			continue
		}

		summary.TotalApps += len(apps)

		// Count versions for each app
		for _, app := range apps {
			versions, err := s.versionRepo.ListByApp(app.ID, 10000, 0)
			if err != nil {
				// Continue with other apps even if one fails
				continue
			}
			summary.TotalVersions += len(versions)
		}
	}

	return summary, nil
}

// buildProjectInventory builds the inventory structure for a single project
func (s *InventoryService) buildProjectInventory(project *database.Project) (*ProjectInventory, error) {
	// Fetch apps for this project
	apps, err := s.appRepo.ListByProject(project.ID, 10000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch apps for project %s: %w", project.Name, err)
	}

	projectInventory := &ProjectInventory{
		Project: *project,
		Apps:    make([]AppInventory, 0, len(apps)),
	}

	// For each app, fetch its versions
	for _, app := range apps {
		versionPtrs, err := s.versionRepo.ListByApp(app.ID, 10000, 0)
		if err != nil {
			// Continue with other apps even if one fails
			// Return empty versions for this app
			projectInventory.Apps = append(projectInventory.Apps, AppInventory{
				App:      *app,
				Versions: []database.Version{},
			})
			continue
		}

		// Convert []*Version to []Version
		versions := make([]database.Version, len(versionPtrs))
		for i, v := range versionPtrs {
			versions[i] = *v
		}

		projectInventory.Apps = append(projectInventory.Apps, AppInventory{
			App:      *app,
			Versions: versions,
		})
	}

	return projectInventory, nil
}

