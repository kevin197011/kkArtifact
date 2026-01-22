// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kk/kkartifact-server/internal/util"
	"github.com/kk/kkartifact-server/internal/version"
)

// findStaticFile tries to find a file in static directory by trying multiple possible paths
func findStaticFile(relativePath string) ([]byte, error) {
	// First, try environment variable for static directory
	if staticDir := os.Getenv("AGENT_STATIC_DIR"); staticDir != "" {
		filePath := filepath.Join(staticDir, relativePath)
		if data, err := os.ReadFile(filePath); err == nil {
			if util.IsDebugMode() {
				log.Printf("Found file via AGENT_STATIC_DIR: %s", filePath)
			}
			return data, nil
		}
		if absPath, err := filepath.Abs(filePath); err == nil {
			if data, err := os.ReadFile(absPath); err == nil {
				if util.IsDebugMode() {
					log.Printf("Found file via AGENT_STATIC_DIR (abs): %s", absPath)
				}
				return data, nil
			}
		}
	}
	
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}
	if util.IsDebugMode() {
		log.Printf("Searching for %s, working directory: %s", relativePath, wd)
	}
	
	// Try multiple possible paths (both relative and absolute)
	possiblePaths := []string{
		// Relative paths (most common)
		filepath.Join("server", "static", relativePath),          // If running from project root
		filepath.Join("static", relativePath),                    // If running from server directory
		filepath.Join("..", "server", "static", relativePath),    // If in subdirectory
		// Absolute paths based on current working directory
		filepath.Join(wd, "server", "static", relativePath),      // Absolute, if in project root
		filepath.Join(wd, "static", relativePath),                // Absolute, if in server directory
		filepath.Join(filepath.Dir(wd), "server", "static", relativePath), // Parent dir
	}
	
	// Try all paths
	for _, path := range possiblePaths {
		// Try as-is first
		if data, err := os.ReadFile(path); err == nil {
			if util.IsDebugMode() {
				log.Printf("Found file at: %s", path)
			}
			return data, nil
		}
		// Try as absolute path
		if absPath, err := filepath.Abs(path); err == nil && absPath != path {
			if data, err := os.ReadFile(absPath); err == nil {
				if util.IsDebugMode() {
					log.Printf("Found file at (abs): %s", absPath)
				}
				return data, nil
			}
		}
	}
	
	if util.IsDebugMode() {
		log.Printf("File not found after trying %d paths", len(possiblePaths))
	}
	return nil, os.ErrNotExist
}

// handleGetAgentVersionInfo returns agent binary version information
// handleGetAgentVersionInfo godoc
// @Summary      Get agent version info
// @Description  Get available agent binaries and version information
// @Tags         downloads
// @Produce      json
// @Success      200  {object}  AgentVersionInfo
// @Failure      500  {object}  ErrorResponse
// @Router       /downloads/agent/version [get]
func (h *Handler) handleGetAgentVersionInfo(c *gin.Context) {
	// Try to read version.json from static/agent directory
	versionData, err := findStaticFile("agent/version.json")
	if err != nil {
		// Log the error for debugging
		wd, _ := os.Getwd()
		if util.IsDebugMode() {
			log.Printf("Failed to find version.json: %v, working directory: %s", err, wd)
		}
		// If version.json doesn't exist, return empty info
		c.JSON(http.StatusOK, gin.H{
			"version": "unknown",
			"build_time": "",
			"binaries": []interface{}{},
		})
		return
	}

	var versionInfo AgentVersionInfo
	if err := json.Unmarshal(versionData, &versionInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse version info"})
		return
	}

	c.JSON(http.StatusOK, versionInfo)
}

// handleDownloadAgent serves agent binary files
// handleDownloadAgent godoc
// @Summary      Download agent binary
// @Description  Download agent binary for specific platform
// @Tags         downloads
// @Param        filename  path  string  true  "Binary filename (e.g., kkartifact-agent-linux-amd64)"
// @Produce      application/octet-stream
// @Success      200  {file}  binary
// @Failure      404  {object}  ErrorResponse
// @Router       /downloads/agent/{filename} [get]
func (h *Handler) handleDownloadAgent(c *gin.Context) {
	filename := c.Param("filename")
	
	// Validate filename to prevent path traversal
	if filename == "" || filepath.Base(filename) != filename {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filename"})
		return
	}

	// Only allow kkartifact-agent files
	matched, err := filepath.Match("kkartifact-agent-*", filename)
	if err != nil || !matched {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filename pattern"})
		return
	}

	// Find the file path using the helper function
	filePath, err := findAgentFilePath(filename)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	// Set headers for file download
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")

	// Serve file
	c.File(filePath)
}

// AgentVersionInfo represents agent version information
type AgentVersionInfo struct {
	Version   string           `json:"version"`
	BuildTime string           `json:"build_time"`
	Binaries  []AgentBinaryInfo `json:"binaries"`
}

// AgentBinaryInfo represents a single agent binary
type AgentBinaryInfo struct {
	Platform string `json:"platform"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	URL      string `json:"url"`
}

// ServerVersionInfo represents server version information
type ServerVersionInfo struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	GitCommit string `json:"git_commit"`
}

// handleGetServerVersion returns server version information
// handleGetServerVersion godoc
// @Summary      Get server version info
// @Description  Get server version, build time, and git commit information
// @Tags         downloads
// @Produce      json
// @Success      200  {object}  ServerVersionInfo
// @Router       /downloads/server/version [get]
func (h *Handler) handleGetServerVersion(c *gin.Context) {
	versionInfo := ServerVersionInfo{
		Version:   version.GetVersion(),
		BuildTime: version.GetBuildTime(),
		GitCommit: version.GetGitCommit(),
	}
	c.JSON(http.StatusOK, versionInfo)
}

// handleDownloadScript serves installation scripts
// handleDownloadScript godoc
// @Summary      Download installation script
// @Description  Download installation script for Unix (install-agent.sh) or Windows (install-agent.ps1)
// @Tags         downloads
// @Param        filename  path  string  true  "Script filename (install-agent.sh or install-agent.ps1)"
// @Produce      text/x-shellscript
// @Produce      application/x-powershell
// @Success      200  {file}  script
// @Failure      404  {object}  ErrorResponse
// @Router       /downloads/scripts/{filename} [get]
func (h *Handler) handleDownloadScript(c *gin.Context) {
	filename := c.Param("filename")
	
	// Validate filename to prevent path traversal
	if filename == "" || filepath.Base(filename) != filename {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filename"})
		return
	}

	// Only allow install-agent scripts
	if filename != "install-agent.sh" && filename != "install-agent.ps1" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid script filename"})
		return
	}

	// Find the script file
	scriptData, err := findStaticFile("scripts/" + filename)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "script not found"})
		return
	}

	// Extract server URL from request
	// Get scheme (http/https)
	scheme := "http"
	if c.GetHeader("X-Forwarded-Proto") == "https" || c.Request.TLS != nil {
		scheme = "https"
	}
	
	// Get host from request
	host := c.GetHeader("Host")
	if host == "" {
		host = c.Request.Host
	}
	if host == "" {
		host = "localhost:8080"
	}
	
	// Construct server URL
	serverURL := scheme + "://" + host

	// Inject server URL into script content
	// Replace __SERVER_URL__ placeholder with actual server URL
	scriptContent := string(scriptData)
	scriptContent = strings.ReplaceAll(scriptContent, "__SERVER_URL__", serverURL)

	// Set appropriate Content-Type based on script type
	if filename == "install-agent.sh" {
		c.Header("Content-Type", "text/x-shellscript")
	} else {
		c.Header("Content-Type", "application/x-powershell")
	}
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")

	// Serve file content with injected server URL
	c.Data(http.StatusOK, c.GetHeader("Content-Type"), []byte(scriptContent))
}
