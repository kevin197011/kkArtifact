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
	"strings"
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
	pullProject    string
	pullApp        string
	pullVersion    string
	pullPath       string
	pullConfig     string
	pullServerURL  string
	pullToken      string
	pullConcurrency int
	pullIgnore     []string
)

func init() {
	rootCmd.AddCommand(pullCmd)
	
	pullCmd.Flags().StringVar(&pullProject, "project", "", "Project name (required)")
	pullCmd.Flags().StringVar(&pullApp, "app", "", "App name (required)")
	pullCmd.Flags().StringVar(&pullVersion, "version", "latest", "Version hash (use 'latest' for latest published version)")
	pullCmd.Flags().StringVar(&pullPath, "path", ".", "Path to local directory")
	pullCmd.Flags().StringVar(&pullConfig, "config", ".kkartifact.yml", "Config file path")
	pullCmd.Flags().StringVar(&pullServerURL, "server-url", "", "Server URL (overrides config file)")
	pullCmd.Flags().StringVar(&pullToken, "token", "", "Authentication token (overrides config file)")
	pullCmd.Flags().IntVar(&pullConcurrency, "concurrency", 0, "Number of concurrent downloads (overrides config file, 0 = use config)")
	pullCmd.Flags().StringArrayVar(&pullIgnore, "ignore", []string{}, "Ignore patterns (can be specified multiple times or comma-separated, merges with config file)")
	
	pullCmd.MarkFlagRequired("project")
	pullCmd.MarkFlagRequired("app")
}

func runPull(cmd *cobra.Command, args []string) error {
	startTime := time.Now()
	
	if pullProject == "" || pullApp == "" {
		return fmt.Errorf("project and app are required")
	}

	// Parse ignore patterns from command-line (support comma-separated values)
	ignorePatterns := make([]string, 0)
	for _, ignoreFlag := range pullIgnore {
		// Split by comma and trim whitespace
		patterns := strings.Split(ignoreFlag, ",")
		for _, pattern := range patterns {
			trimmed := strings.TrimSpace(pattern)
			if trimmed != "" {
				ignorePatterns = append(ignorePatterns, trimmed)
			}
		}
	}

	// Prepare command-line overrides
	overrides := &config.Overrides{
		ServerURL:   pullServerURL,
		Token:        pullToken,
		Concurrency: pullConcurrency,
		Ignore:      ignorePatterns,
	}
	if len(ignorePatterns) == 0 {
		overrides.Ignore = nil // Don't override if no ignore patterns provided
	}
	if pullConcurrency == 0 {
		overrides.Concurrency = 0 // 0 means not set
	}

	// Load config with overrides
	cfg, err := config.Load(pullConfig, overrides)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate token is set
	if cfg.Token == "" {
		return fmt.Errorf("token is required but not found in config. Please check:\n  - Global config: /etc/kkArtifact/config.yml\n  - Local config: %s\n  - Or use --token flag", pullConfig)
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

	// Create API client
	apiClient := client.New(cfg.ServerURL, cfg.Token)

	// Handle "latest" version
	actualVersion := pullVersion
	if pullVersion == "" || pullVersion == "latest" {
		fmt.Printf("Fetching latest published version for %s/%s...\n", pullProject, pullApp)
		latestResp, err := apiClient.GetLatestVersion(pullProject, pullApp)
		if err != nil {
			return fmt.Errorf("failed to get latest version: %w", err)
		}
		actualVersion = latestResp.Version
		fmt.Printf("Latest published version: %s\n", actualVersion)
	}

	fmt.Printf("Pulling artifacts from %s/%s:%s to %s\n", pullProject, pullApp, actualVersion, absPath)

	// Get manifest
	fmt.Println("Fetching manifest...")
	manifestData, err := apiClient.GetManifest(pullProject, pullApp, actualVersion)
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
	
	// Create progress bar
	progressBar := NewProgressBar(len(manifest.Files))
	
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
				exists, matches, _, err := client.CheckFileExistsAndMatches(task.localPath, task.expectedHash)
				if err != nil {
					errors <- fmt.Errorf("failed to check file %s: %w", task.filePath, err)
					return
				}
				
				if exists && matches {
					// Skip file - update progress bar
					progressBar.Update(1)
					continue
				}
				
				// Download file (with resume support if partial file exists)
				if err := apiClient.DownloadFile(pullProject, pullApp, actualVersion, task.filePath, task.localPath, task.expectedHash, task.expectedSize); err != nil {
					errors <- fmt.Errorf("failed to download file %s: %w", task.filePath, err)
					return
				}
				
				// Update progress bar
				progressBar.Update(1)
			}
		}()
	}
	
	// Wait for all workers to finish
	wg.Wait()
	
	// Finish progress bar
	progressBar.Finish()
	
	close(errors)
	
	// Check for errors
	for err := range errors {
		if err != nil {
			return err
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("Successfully pulled %s/%s:%s\n", pullProject, pullApp, actualVersion)
	fmt.Printf("Total time: %v\n", duration.Round(time.Second))
	return nil
}
