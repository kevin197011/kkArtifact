// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"os"
	"path/filepath"
)

// findAgentFilePath finds the full path to an agent file
func findAgentFilePath(filename string) (string, error) {
	// First, try environment variable for static directory
	if staticDir := os.Getenv("AGENT_STATIC_DIR"); staticDir != "" {
		filePath := filepath.Join(staticDir, "agent", filename)
		if _, err := os.Stat(filePath); err == nil {
			return filePath, nil
		}
		if absPath, err := filepath.Abs(filePath); err == nil {
			if _, err := os.Stat(absPath); err == nil {
				return absPath, nil
			}
		}
	}
	
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}
	
	// Try multiple possible paths
	// Priority: Docker container paths first, then development paths
	possiblePaths := []string{
		// Docker container paths (WORKDIR /app)
		filepath.Join("/app", "static", "agent", filename),
		filepath.Join("static", "agent", filename), // Relative to working directory
		// Development paths
		filepath.Join("server", "static", "agent", filename),
		filepath.Join(wd, "static", "agent", filename),
		filepath.Join(wd, "server", "static", "agent", filename),
	}
	
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
		if absPath, err := filepath.Abs(path); err == nil {
			if _, err := os.Stat(absPath); err == nil {
				return absPath, nil
			}
		}
	}
	
	return "", os.ErrNotExist
}

