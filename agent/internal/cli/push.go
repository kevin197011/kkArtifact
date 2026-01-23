// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/kk/kkartifact-agent/internal/client"
	"github.com/kk/kkartifact-agent/internal/config"
	"github.com/kk/kkartifact-agent/internal/manifest"
)

var pushCmd = &cobra.Command{
	Use:          "push [flags]",
	Short:        "Push artifacts to the server",
	Long:         "Push artifacts from a local directory to the kkArtifact server",
	SilenceUsage: true, // Don't show usage on errors
	RunE:         runPush,
}

var (
	pushProject    string
	pushApp        string
	pushVersion    string
	pushPath       string
	pushConfig     string
	pushServerURL  string
	pushToken      string
	pushConcurrency int
	pushIgnore     []string
)

func init() {
	rootCmd.AddCommand(pushCmd)
	
	pushCmd.Flags().StringVar(&pushProject, "project", "", "Project name (required)")
	pushCmd.Flags().StringVar(&pushApp, "app", "", "App name (required)")
	pushCmd.Flags().StringVar(&pushVersion, "version", "", "Version hash (required)")
	pushCmd.Flags().StringVar(&pushPath, "path", ".", "Path to local directory")
	pushCmd.Flags().StringVar(&pushConfig, "config", ".kkartifact.yml", "Config file path")
	pushCmd.Flags().StringVar(&pushServerURL, "server-url", "", "Server URL (overrides config file)")
	pushCmd.Flags().StringVar(&pushToken, "token", "", "Authentication token (overrides config file)")
	pushCmd.Flags().IntVar(&pushConcurrency, "concurrency", 0, "Number of concurrent uploads (overrides config file, 0 = use config)")
	pushCmd.Flags().StringArrayVar(&pushIgnore, "ignore", []string{}, "Ignore patterns (can be specified multiple times or comma-separated, merges with config file)")
	
	pushCmd.MarkFlagRequired("project")
	pushCmd.MarkFlagRequired("app")
	pushCmd.MarkFlagRequired("version")
}

func runPush(cmd *cobra.Command, args []string) error {
	startTime := time.Now()
	
	if pushProject == "" || pushApp == "" || pushVersion == "" {
		return fmt.Errorf("project, app, and version are required")
	}

	// Parse ignore patterns from command-line (support comma-separated values)
	ignorePatterns := make([]string, 0)
	for _, ignoreFlag := range pushIgnore {
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
		ServerURL:   pushServerURL,
		Token:        pushToken,
		Concurrency: pushConcurrency,
		Ignore:      ignorePatterns,
	}
	if len(ignorePatterns) == 0 {
		overrides.Ignore = nil // Don't override if no ignore patterns provided
	}
	if pushConcurrency == 0 {
		overrides.Concurrency = 0 // 0 means not set
	}

	// Load config with overrides
	cfg, err := config.Load(pushConfig, overrides)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate token is set
	if cfg.Token == "" {
		return fmt.Errorf("token is required but not found in config. Please check:\n  - Global config: /etc/kkArtifact/config.yml\n  - Local config: %s\n  - Or use --token flag", pushConfig)
	}

	// Validate token format before creating client
	if err := config.ValidateTokenFormat(cfg.Token); err != nil {
		return fmt.Errorf("token validation failed: %w\nToken preview: %s\nConfig file: %s", err, config.MaskToken(cfg.Token), pushConfig)
	}

	// Use config values if not provided via flags
	if pushProject == "" {
		pushProject = cfg.Project
	}
	if pushApp == "" {
		pushApp = cfg.App
	}

	// Resolve absolute path
	absPath, err := filepath.Abs(pushPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check if path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", absPath)
	}

	fmt.Printf("Generating manifest for %s/%s:%s from %s\n", pushProject, pushApp, pushVersion, absPath)

	// Generate manifest
	m, err := manifest.Generate(pushProject, pushApp, pushVersion, absPath, cfg.Ignore)
	if err != nil {
		return fmt.Errorf("failed to generate manifest: %w", err)
	}

	fmt.Printf("Found %d files\n", len(m.Files))

	// Create API client with validation
	apiClient, err := client.New(cfg.ServerURL, cfg.Token)
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	// Initialize upload
	fmt.Println("Initializing upload...")
	uploadResp, err := apiClient.InitUpload(pushProject, pushApp, pushVersion, len(m.Files))
	if err != nil {
		return fmt.Errorf("failed to initialize upload: %w", err)
	}
	fmt.Printf("Upload ID: %s\n", uploadResp.UploadID)

	// Upload files concurrently
	fmt.Printf("Uploading %d files with concurrency: %d\n", len(m.Files), cfg.Concurrency)
	
	// Create progress bar
	progressBar := NewProgressBar(len(m.Files))
	
	type uploadTask struct {
		index    int
		file     manifest.ManifestFile
		localPath string
	}
	
	tasks := make(chan uploadTask, len(m.Files))
	errors := make(chan error, len(m.Files))
	
	// Populate tasks
	for i, file := range m.Files {
		localPath := filepath.Join(absPath, file.Path)
		tasks <- uploadTask{
			index:     i,
			file:      file,
			localPath: localPath,
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
				// Calculate local file hash to verify it matches manifest
				// Note: We always upload since server handles overwrite, but we could skip if hash matches
				// For now, we upload all files as the server already handles overwrite in handleInitUpload
				
				if err := apiClient.UploadFile(pushProject, pushApp, pushVersion, task.file.Path, task.localPath); err != nil {
					errors <- fmt.Errorf("failed to upload file %s: %w", task.file.Path, err)
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

	// Finish upload
	fmt.Println("Finalizing upload...")
	finishReq := map[string]interface{}{
		"project":  pushProject,
		"app":      pushApp,
		"version":  pushVersion,
		"manifest": m,
	}

	if err := apiClient.FinishUpload(finishReq); err != nil {
		return fmt.Errorf("failed to finish upload: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("Successfully pushed %s/%s:%s\n", pushProject, pushApp, pushVersion)
	fmt.Printf("Total time: %v\n", duration.Round(time.Second))
	return nil
}
