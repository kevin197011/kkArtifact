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

// handleSyncStorage syncs database records from storage files
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

	projectCount := 0
	appCount := 0
	versionCount := 0
	processedProjects := make(map[string]bool)
	processedApps := make(map[string]bool)

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
				// No manifest, skip
				return nil
			}

			project, err := h.projectRepo.CreateOrGet(projectName)
			if err != nil {
				log.Printf("Warning: failed to create/get project %s: %v", projectName, err)
				return nil
			}

			app, err := h.appRepo.CreateOrGet(project.ID, appName)
			if err != nil {
				log.Printf("Warning: failed to create/get app %s/%s: %v", projectName, appName, err)
				return nil
			}

			_, err = h.versionRepo.Create(app.ID, versionHash)
			if err != nil {
				// Version might already exist, ignore error
			} else {
				versionCount++
				log.Printf("    Added version: %s/%s/%s", projectName, appName, versionHash)
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("Error during sync: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "sync completed",
		"projects": projectCount,
		"apps":     appCount,
		"versions": versionCount,
	})
}

