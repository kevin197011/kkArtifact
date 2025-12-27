// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ManifestResponse represents a manifest in API response
type ManifestResponse struct {
	Project   string                 `json:"project"`
	App       string                 `json:"app"`
	Version   string                 `json:"version"`
	Hash      string                 `json:"hash"`      // Same as version
	GitCommit string                 `json:"git_commit,omitempty"`
	BuildTime string                 `json:"build_time"`
	Builder   string                 `json:"builder"`
	Files     []ManifestFileResponse `json:"files"`
}

// ManifestFileResponse represents a file in manifest response
type ManifestFileResponse struct {
	Path   string `json:"path"`
	Hash   string `json:"hash"`   // SHA256 hash
	Size   int64  `json:"size"`
}

func (h *Handler) handleGetManifest(c *gin.Context) {
	project := c.Param("project")
	app := c.Param("app")
	hash := c.Param("hash")

	manifest, err := h.artifactManager.GetManifest(c.Request.Context(), project, app, hash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	files := make([]ManifestFileResponse, len(manifest.Files))
	for i, f := range manifest.Files {
		files[i] = ManifestFileResponse{
			Path: f.Path,
			Hash: f.SHA256, // Convert SHA256 to hash
			Size: f.Size,
		}
	}

	response := ManifestResponse{
		Project:   manifest.Project,
		App:       manifest.App,
		Version:   manifest.Version,
		Hash:      manifest.Version, // Use version as hash for compatibility
		GitCommit: manifest.GitCommit,
		BuildTime: manifest.BuildTime,
		Builder:   manifest.Builder,
		Files:     files,
	}

	c.JSON(http.StatusOK, response)
}

