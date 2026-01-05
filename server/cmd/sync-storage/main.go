// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kk/kkartifact-server/internal/config"
	"github.com/kk/kkartifact-server/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	projectRepo := database.NewProjectRepository(db)
	appRepo := database.NewAppRepository(db)
	versionRepo := database.NewVersionRepository(db)

	// Scan storage directory structure using filepath.Walk
	// Storage structure: {project}/{app}/{version}/meta.yaml
	basePath := cfg.Storage.LocalPath
	log.Printf("Scanning storage at: %s", basePath)

	projectCount := 0
	appCount := 0
	versionCount := 0
	processedProjects := make(map[string]bool)
	processedApps := make(map[string]bool)

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
				_, err := projectRepo.CreateOrGet(projectName)
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
				project, err := projectRepo.CreateOrGet(projectName)
				if err != nil {
					log.Printf("Warning: failed to create/get project %s: %v", projectName, err)
					return nil
				}

				_, err = appRepo.CreateOrGet(project.ID, appName)
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

			project, err := projectRepo.CreateOrGet(projectName)
			if err != nil {
				log.Printf("Warning: failed to create/get project %s: %v", projectName, err)
				return nil
			}

			app, err := appRepo.CreateOrGet(project.ID, appName)
			if err != nil {
				log.Printf("Warning: failed to create/get app %s/%s: %v", projectName, appName, err)
				return nil
			}

			_, err = versionRepo.Create(app.ID, versionHash)
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
		log.Fatalf("Failed to scan storage: %v", err)
	}

	log.Printf("========================================")
	log.Printf("Sync completed:")
	log.Printf("  Projects: %d", projectCount)
	log.Printf("  Apps: %d", appCount)
	log.Printf("  Versions: %d", versionCount)
	log.Printf("========================================")
}
