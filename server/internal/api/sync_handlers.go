// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// SyncStorageResponse represents the response for sync storage operation
type SyncStorageResponse struct {
	Message  string `json:"message"`
	Projects int    `json:"projects"`
	Apps     int    `json:"apps"`
	Versions int    `json:"versions"`
}

// handleSyncStorage godoc
// @Summary      Sync storage
// @Description  Synchronize database with storage (rebuild database records from storage, remove orphaned records)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  SyncStorageResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     Bearer
// @Router       /sync-storage [post]
func (h *Handler) handleSyncStorage(c *gin.Context) {
	// Only allow sync for local storage
	if h.storage == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "storage not initialized"})
		return
	}

	// Get storage base path from config
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType != "local" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sync is only supported for local storage"})
		return
	}

	basePath := os.Getenv("STORAGE_LOCAL_PATH")
	if basePath == "" {
		basePath = "/repos"
	}

	log.Printf("Starting storage sync from: %s", basePath)

	// Maps to track what exists in storage
	storageProjects := make(map[string]bool)                         // project name -> exists
	storageApps := make(map[string]bool)                              // "project/app" -> exists
	storageVersions := make(map[string]map[string]map[string]bool)    // project -> app -> version hash -> exists

	projectCount := 0
	appCount := 0
	versionCount := 0
	processedProjects := make(map[string]bool)
	processedApps := make(map[string]bool)

	// First pass: scan storage and collect what exists
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		// Get relative path from base
		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return nil
		}

		// Skip root directory
		if relPath == "." {
			return nil
		}

		parts := strings.Split(relPath, string(filepath.Separator))

		// Structure: project/app/version
		if len(parts) == 1 {
			// Project level
			projectName := parts[0]
			storageProjects[projectName] = true
		} else if len(parts) == 2 {
			// App level
			projectName := parts[0]
			appName := parts[1]
			key := projectName + "/" + appName
			storageApps[key] = true
		} else if len(parts) == 3 {
			// Version level - check if meta.yaml exists
			projectName := parts[0]
			appName := parts[1]
			versionHash := parts[2]

			manifestPath := filepath.Join(path, "meta.yaml")
			if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
				// No manifest, skip this version directory
				// Skip deeper traversal into version directory to avoid processing subdirectories
				return filepath.SkipDir
			}

			// Track this version exists (has meta.yaml)
			if storageVersions[projectName] == nil {
				storageVersions[projectName] = make(map[string]map[string]bool)
			}
			if storageVersions[projectName][appName] == nil {
				storageVersions[projectName][appName] = make(map[string]bool)
			}
			storageVersions[projectName][appName][versionHash] = true
			// Skip deeper traversal into version directory to avoid processing subdirectories
			return filepath.SkipDir
		} else if len(parts) > 3 {
			// Deeper than version level, skip to avoid processing subdirectories
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		log.Printf("Error during storage scan: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Second pass: add/update records that exist in storage
	err = filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		// Get relative path from base
		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return nil
		}

		// Skip root directory
		if relPath == "." {
			return nil
		}

		parts := strings.Split(relPath, string(filepath.Separator))

		// Structure: project/app/version
		if len(parts) == 1 {
			// Project level
			projectName := parts[0]
			if !processedProjects[projectName] {
				_, err := h.projectRepo.CreateOrGet(projectName)
				if err != nil {
					log.Printf("Warning: failed to create/get project %s: %v", projectName, err)
				} else {
					projectCount++
					processedProjects[projectName] = true
					log.Printf("Processing project: %s", projectName)
				}
			}
		} else if len(parts) == 2 {
			// App level
			projectName := parts[0]
			appName := parts[1]
			key := projectName + "/" + appName
			if !processedApps[key] {
				project, err := h.projectRepo.CreateOrGet(projectName)
				if err != nil {
					log.Printf("Warning: failed to create/get project %s: %v", projectName, err)
					return nil
				}

				_, err = h.appRepo.CreateOrGet(project.ID, appName)
				if err != nil {
					log.Printf("Warning: failed to create/get app %s/%s: %v", projectName, appName, err)
				} else {
					appCount++
					processedApps[key] = true
					log.Printf("  Processing app: %s/%s", projectName, appName)
				}
			}
		} else if len(parts) == 3 {
			// Version level - check if meta.yaml exists
			projectName := parts[0]
			appName := parts[1]
			versionHash := parts[2]

			manifestPath := filepath.Join(path, "meta.yaml")
			if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
				// No manifest, skip this version (only versions with meta.yaml are considered complete)
				// Skip deeper traversal into version directory to avoid processing subdirectories
				return filepath.SkipDir
			}

			project, err := h.projectRepo.CreateOrGet(projectName)
			if err != nil {
				log.Printf("Warning: failed to create/get project %s: %v", projectName, err)
				return filepath.SkipDir
			}

			app, err := h.appRepo.CreateOrGet(project.ID, appName)
			if err != nil {
				log.Printf("Warning: failed to create/get app %s/%s: %v", projectName, appName, err)
				return filepath.SkipDir
			}

			_, err = h.versionRepo.Create(app.ID, versionHash)
			if err != nil {
				// Version might already exist, ignore error
			} else {
				versionCount++
				log.Printf("    Added version: %s/%s/%s", projectName, appName, versionHash)
			}
			
			// Skip deeper traversal into version directory to avoid processing subdirectories
			return filepath.SkipDir
		} else if len(parts) > 3 {
			// Deeper than version level, skip to avoid processing subdirectories
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		log.Printf("Error during sync: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Third pass: delete records that don't exist in storage
	// Get all projects from database
	allProjects, err := h.projectRepo.List(10000, 0)
	if err != nil {
		log.Printf("Warning: failed to list projects: %v", err)
	} else {
		deletedProjectCount := 0
		for _, project := range allProjects {
			if !storageProjects[project.Name] {
				// Project doesn't exist in storage, delete it (cascade will delete apps and versions)
				query := `DELETE FROM projects WHERE id = $1`
				_, err := h.db.Exec(query, project.ID)
				if err != nil {
					log.Printf("Warning: failed to delete project %s: %v", project.Name, err)
				} else {
					deletedProjectCount++
					log.Printf("Deleted project (not in storage): %s", project.Name)
				}
			}
		}
		if deletedProjectCount > 0 {
			log.Printf("Deleted %d projects that don't exist in storage", deletedProjectCount)
		}
	}

	// Get all apps from database and check against storage
	allProjects, err = h.projectRepo.List(10000, 0)
	if err == nil {
		deletedAppCount := 0
		for _, project := range allProjects {
			apps, err := h.appRepo.ListByProject(project.ID, 10000, 0)
			if err != nil {
				log.Printf("Warning: failed to list apps for project %s: %v", project.Name, err)
				continue
			}
			for _, app := range apps {
				appKey := project.Name + "/" + app.Name
				if !storageApps[appKey] {
					// App doesn't exist in storage, delete it (cascade will delete versions)
					query := `DELETE FROM apps WHERE id = $1`
					_, err := h.db.Exec(query, app.ID)
					if err != nil {
						log.Printf("Warning: failed to delete app %s/%s: %v", project.Name, app.Name, err)
					} else {
						deletedAppCount++
						log.Printf("Deleted app (not in storage): %s/%s", project.Name, app.Name)
					}
				}
			}
		}
		if deletedAppCount > 0 {
			log.Printf("Deleted %d apps that don't exist in storage", deletedAppCount)
		}
	}

	// Get all versions from database and check against storage
	allProjects, err = h.projectRepo.List(10000, 0)
	if err == nil {
		deletedVersionCount := 0
		for _, project := range allProjects {
			apps, err := h.appRepo.ListByProject(project.ID, 10000, 0)
			if err != nil {
				continue
			}
			for _, app := range apps {
				versions, err := h.versionRepo.ListByApp(app.ID, 10000, 0)
				if err != nil {
					continue
				}
				for _, version := range versions {
					// Check if this version exists in storage
					exists := false
					if storageVersions[project.Name] != nil && storageVersions[project.Name][app.Name] != nil {
						exists = storageVersions[project.Name][app.Name][version.Hash]
					}
					if !exists {
						// Version doesn't exist in storage, delete it
						err := h.versionRepo.Delete(app.ID, version.Hash)
						if err != nil {
							log.Printf("Warning: failed to delete version %s/%s/%s: %v", project.Name, app.Name, version.Hash, err)
						} else {
							deletedVersionCount++
							log.Printf("Deleted version (not in storage): %s/%s/%s", project.Name, app.Name, version.Hash)
						}
					}
				}
			}
		}
		if deletedVersionCount > 0 {
			log.Printf("Deleted %d versions that don't exist in storage", deletedVersionCount)
		}
	}

	c.JSON(http.StatusOK, SyncStorageResponse{
		Message:  "sync completed",
		Projects: projectCount,
		Apps:     appCount,
		Versions: versionCount,
	})
}

