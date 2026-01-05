// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/kk/kkartifact-agent/internal/client"
	"github.com/kk/kkartifact-agent/internal/config"
)

var updateCmd = &cobra.Command{
	Use:          "update [flags]",
	Short:        "Update agent to the latest version",
	Long:         "Download and install the latest version of kkartifact-agent from the server",
	SilenceUsage: true,
	RunE:         runUpdate,
}

var (
	updateConfig string
	updateForce  bool
)

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVar(&updateConfig, "config", ".kkartifact.yml", "Config file path")
	updateCmd.Flags().BoolVar(&updateForce, "force", false, "Force update even if already on latest version")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	// Load config to get server URL
	cfg, err := config.Load(updateConfig)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create API client (no token needed for public endpoints)
	apiClient := client.New(cfg.ServerURL, "")

	// Get current binary path
	currentBinary, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current binary path: %w", err)
	}
	absBinary, err := filepath.Abs(currentBinary)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	fmt.Printf("Current agent binary: %s\n", absBinary)

	// Get version info from server
	fmt.Println("Fetching latest version information...")
	versionInfo, err := apiClient.GetAgentVersionInfo()
	if err != nil {
		return fmt.Errorf("failed to get version info: %w", err)
	}

	if versionInfo.Version == "" || len(versionInfo.Binaries) == 0 {
		return fmt.Errorf("invalid version info from server")
	}

	fmt.Printf("Latest version available: %s (built at %s)\n", versionInfo.Version, versionInfo.BuildTime)

	// Get current version (try to get from binary if possible)
	currentVersion := getCurrentVersion()
	if currentVersion != "" {
		fmt.Printf("Current version: %s\n", currentVersion)
		if currentVersion == versionInfo.Version && !updateForce {
			fmt.Println("Already on latest version. Use --force to update anyway.")
			return nil
		}
	}

	// Determine platform
	platform := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Platform: %s\n", platform)

	// Find matching binary
	var targetBinary *client.AgentBinaryInfo
	for i := range versionInfo.Binaries {
		if versionInfo.Binaries[i].Platform == platform {
			targetBinary = &versionInfo.Binaries[i]
			break
		}
	}

	if targetBinary == nil {
		return fmt.Errorf("no binary available for platform %s", platform)
	}

	fmt.Printf("Target binary: %s (size: %d bytes)\n", targetBinary.Filename, targetBinary.Size)

	// Create temporary file for download
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, targetBinary.Filename+".tmp")
	defer os.Remove(tempFile) // Clean up on exit

	// Download binary
	fmt.Printf("Downloading %s...\n", targetBinary.Filename)

	// Use the client's download method
	if err := apiClient.DownloadAgentBinary(targetBinary.Filename, tempFile); err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}

	// Verify file size
	fileInfo, err := os.Stat(tempFile)
	if err != nil {
		return fmt.Errorf("failed to stat downloaded file: %w", err)
	}
	if fileInfo.Size() != targetBinary.Size {
		return fmt.Errorf("downloaded file size mismatch: expected %d, got %d", targetBinary.Size, fileInfo.Size())
	}

	fmt.Println("Download completed successfully")

	// Replace current binary
	fmt.Printf("Replacing binary at %s...\n", absBinary)

	// On Windows, we need to handle file locking differently
	if runtime.GOOS == "windows" {
		// On Windows, we can't replace the running executable directly
		// We'll create a batch script to do it after exit
		batchScript := filepath.Join(tempDir, "kkartifact-update.bat")
		batchContent := fmt.Sprintf(`@echo off
timeout /t 2 /nobreak >nul
copy /Y "%s" "%s"
del "%%~f0"
`, tempFile, absBinary)

		if err := os.WriteFile(batchScript, []byte(batchContent), 0755); err != nil {
			return fmt.Errorf("failed to create update script: %w", err)
		}

		// Execute batch script in background
		cmd := exec.Command("cmd", "/C", "start", "/B", batchScript)
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start update script: %w", err)
		}

		fmt.Println("Update will complete after agent exits. Please restart the agent.")
		return nil
	}

	// On Unix systems, we can replace the file directly
	// First, make sure the temp file is executable
	if err := os.Chmod(tempFile, 0755); err != nil {
		return fmt.Errorf("failed to make file executable: %w", err)
	}

	// Replace the binary atomically
	backupPath := absBinary + ".backup"
	if err := os.Rename(absBinary, backupPath); err != nil {
		return fmt.Errorf("failed to backup current binary: %w", err)
	}

	if err := os.Rename(tempFile, absBinary); err != nil {
		// Restore backup on error
		os.Rename(backupPath, absBinary)
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	// Remove backup
	os.Remove(backupPath)

	fmt.Printf("Successfully updated to version %s\n", versionInfo.Version)
	fmt.Println("Please restart the agent to use the new version.")

	return nil
}

// getCurrentVersion tries to get the current version of the running binary
// This is a placeholder - in a real implementation, you might embed version info at build time
func getCurrentVersion() string {
	// Try to get version from binary metadata or build tags
	// For now, return empty string to always check for updates
	// In production, you could:
	// 1. Embed version at build time using -ldflags
	// 2. Read from a version file
	// 3. Query the binary itself
	return ""
}

