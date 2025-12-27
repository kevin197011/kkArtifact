// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/kk/kkartifact-agent/internal/client"
	"github.com/kk/kkartifact-agent/internal/config"
)

var pullCmd = &cobra.Command{
	Use:          "pull [flags]",
	Short:        "Pull artifacts from the server",
	Long:         "Pull artifacts from the kkArtifact server to a local directory",
	SilenceUsage: true, // Don't show usage on errors
	RunE:         runPull,
}

var (
	pullProject string
	pullApp     string
	pullVersion string
	pullPath    string
	pullConfig  string
)

func init() {
	rootCmd.AddCommand(pullCmd)
	
	pullCmd.Flags().StringVar(&pullProject, "project", "", "Project name (required)")
	pullCmd.Flags().StringVar(&pullApp, "app", "", "App name (required)")
	pullCmd.Flags().StringVar(&pullVersion, "version", "", "Version hash (required)")
	pullCmd.Flags().StringVar(&pullPath, "path", ".", "Path to local directory")
	pullCmd.Flags().StringVar(&pullConfig, "config", ".kkartifact.yml", "Config file path")
	
	pullCmd.MarkFlagRequired("project")
	pullCmd.MarkFlagRequired("app")
	pullCmd.MarkFlagRequired("version")
}

func runPull(cmd *cobra.Command, args []string) error {
	startTime := time.Now()
	
	if pullProject == "" || pullApp == "" || pullVersion == "" {
		return fmt.Errorf("project, app, and version are required")
	}

	// Load config
	cfg, err := config.Load(pullConfig)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Use config values if not provided via flags
	if pullProject == "" {
		pullProject = cfg.Project
	}
	if pullApp == "" {
		pullApp = cfg.App
	}

	// Resolve absolute path
	absPath, err := filepath.Abs(pullPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("Pulling artifacts from %s/%s:%s to %s\n", pullProject, pullApp, pullVersion, absPath)

	// Create API client
	apiClient := client.New(cfg.ServerURL, cfg.Token)

	// Get manifest
	fmt.Println("Fetching manifest...")
	manifestData, err := apiClient.GetManifest(pullProject, pullApp, pullVersion)
	if err != nil {
		return fmt.Errorf("failed to get manifest: %w", err)
	}

	// Parse manifest (convert to map for now)
	manifestBytes, err := json.Marshal(manifestData)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	var manifest struct {
		Files []struct {
			Path   string `json:"path"`
			SHA256 string `json:"hash"` // Note: server returns "hash" not "sha256"
			Size   int64  `json:"size"`
		} `json:"files"`
	}

	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	fmt.Printf("Found %d files in manifest\n", len(manifest.Files))

	// Download files concurrently with resume support
	fmt.Printf("Downloading %d files with concurrency: %d (resume enabled)\n", len(manifest.Files), cfg.Concurrency)
	
	type downloadTask struct {
		index        int
		filePath     string
		localPath    string
		expectedHash string
		expectedSize int64
	}
	
	tasks := make(chan downloadTask, len(manifest.Files))
	errors := make(chan error, len(manifest.Files))
	
	// Populate tasks
	for i, file := range manifest.Files {
		localPath := filepath.Join(absPath, file.Path)
		tasks <- downloadTask{
			index:        i,
			filePath:     file.Path,
			localPath:    localPath,
			expectedHash: file.SHA256,
			expectedSize: file.Size,
		}
	}
	close(tasks)
	
	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				// Check if file needs download
				exists, matches, size, err := client.CheckFileExistsAndMatches(task.localPath, task.expectedHash)
				if err != nil {
					errors <- fmt.Errorf("failed to check file %s: %w", task.filePath, err)
					return
				}
				
				if exists && matches {
					fmt.Printf("[%d/%d] Skipping %s (already exists and matches)\n", task.index+1, len(manifest.Files), task.filePath)
					continue
				}
				
				if exists && size > 0 && task.expectedSize > 0 && size < task.expectedSize {
					percentage := float64(size) * 100 / float64(task.expectedSize)
					fmt.Printf("[%d/%d] Resuming %s (%.1f%% complete)\n", task.index+1, len(manifest.Files), task.filePath, percentage)
				} else {
					fmt.Printf("[%d/%d] Downloading %s...\n", task.index+1, len(manifest.Files), task.filePath)
				}
				
				if err := apiClient.DownloadFile(pullProject, pullApp, pullVersion, task.filePath, task.localPath, task.expectedHash, task.expectedSize); err != nil {
					errors <- fmt.Errorf("failed to download file %s: %w", task.filePath, err)
					return
				}
			}
		}()
	}
	
	// Wait for all workers to finish
	wg.Wait()
	close(errors)
	
	// Check for errors
	for err := range errors {
		if err != nil {
			return err
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("Successfully pulled %s/%s:%s\n", pullProject, pullApp, pullVersion)
	fmt.Printf("Total time: %v\n", duration.Round(time.Second))
	return nil
}
