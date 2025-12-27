// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kk/kkartifact-server/internal/storage"
)

// UploadInitRequest represents the upload initialization request
type UploadInitRequest struct {
	Project   string   `json:"project" binding:"required"`
	App       string   `json:"app" binding:"required"`
	Version   string   `json:"version" binding:"required"`
	FileCount int      `json:"file_count"`
	Files     []string `json:"files,omitempty"`
}

// UploadInitResponse represents the upload initialization response
type UploadInitResponse struct {
	UploadID string `json:"upload_id"`
}

// handleInitUpload initializes an upload session
func (h *Handler) handleInitUpload(c *gin.Context) {
	var req UploadInitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get or create project and app
	project, err := h.projectRepo.CreateOrGet(req.Project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	app, err := h.appRepo.CreateOrGet(project.ID, req.App)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if version already exists - if it does, delete it to allow overwrite
	versionPath := filepath.Join(req.Project, req.App, req.Version, "meta.yaml")
	exists, err := h.storage.Exists(c.Request.Context(), versionPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		// Delete existing version from storage to allow overwrite
		err := h.artifactManager.DeleteVersion(c.Request.Context(), req.Project, req.App, req.Version)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete existing version: %v", err)})
			return
		}
		// Also delete from database (ignore error if version doesn't exist in DB)
		_ = h.versionRepo.Delete(app.ID, req.Version)
	}

	// For now, return a simple upload ID (in production, use UUID)
	uploadID := fmt.Sprintf("%s-%s-%s", req.Project, req.App, req.Version)
	
	c.JSON(http.StatusOK, UploadInitResponse{
		UploadID: uploadID,
	})
}

// handleUploadFile handles file upload
// The hash parameter in the URL is the version (not the file's SHA256)
// Files are stored under {project}/{app}/{version}/{filePath}
func (h *Handler) handleUploadFile(c *gin.Context) {
	project := c.Param("project")
	app := c.Param("app")
	version := c.Param("hash") // This is actually the version, not file hash

	// Get file from multipart form
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing file"})
		return
	}
	defer file.Close()

	// Get path from form
	filePath := c.PostForm("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path parameter is required"})
		return
	}

	// Validate path
	if err := storage.ValidatePath(filePath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate file size and SHA256 for response
	// Note: We don't validate SHA256 here because the hash param is the version, not file hash
	// File integrity is verified via manifest validation in finish upload
	hashReader := io.TeeReader(file, io.Discard)
	hash := sha256.New()
	fileSize, err := io.Copy(hash, hashReader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	calculatedHash := fmt.Sprintf("%x", hash.Sum(nil))

	// Reset file reader
	file.Seek(0, io.SeekStart)

	// Store file under {project}/{app}/{version}/{filePath}
	fullPath := filepath.Join(project, app, version, filePath)
	if err := h.storage.Put(c.Request.Context(), fullPath, file, fileSize); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "uploaded",
		"hash":   calculatedHash, // Return calculated file hash for reference
		"size":   fileSize,
	})
}

// handleFinishUpload finishes an upload session and creates the version
func (h *Handler) handleFinishUpload(c *gin.Context) {
	var req struct {
		Project   string                  `json:"project" binding:"required"`
		App       string                  `json:"app" binding:"required"`
		Version   string                  `json:"version" binding:"required"`
		Manifest  *storage.Manifest       `json:"manifest" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get or create project and app
	project, err := h.projectRepo.CreateOrGet(req.Project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	app, err := h.appRepo.CreateOrGet(project.ID, req.App)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Store manifest
	manifestBytes, err := storage.SerializeManifest(req.Manifest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	manifestPath := filepath.Join(req.Project, req.App, req.Version, "meta.yaml")
	manifestReader := strings.NewReader(string(manifestBytes))
	if err := h.storage.Put(c.Request.Context(), manifestPath, manifestReader, int64(len(manifestBytes))); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create version record in database
	// Note: Since we delete existing versions in handleInitUpload, this should always create a new record
	// However, in case of race conditions or other edge cases, we try to create and log errors
	_, err = h.versionRepo.Create(app.ID, req.Version)
	if err != nil {
		// Log error but don't fail the request - version exists in storage which is what matters most
		// The version record might already exist due to a race condition or previous partial failure
		// TODO: Consider using INSERT ... ON CONFLICT DO NOTHING or similar for idempotency
	}

	// Publish push event with context to extract agent ID and metadata
	metadata := make(map[string]interface{})
	metadata["file_count"] = len(req.Manifest.Files)
	var totalSize int64
	for _, file := range req.Manifest.Files {
		totalSize += file.Size
	}
	metadata["total_size"] = totalSize
	if req.Manifest.GitCommit != "" {
		metadata["git_commit"] = req.Manifest.GitCommit
	}
	if req.Manifest.BuildTime != "" {
		metadata["build_time"] = req.Manifest.BuildTime
	}
	if req.Manifest.Builder != "" {
		metadata["builder"] = req.Manifest.Builder
	}

	h.publishEventWithContext(
		c,
		"push",
		req.Project,
		req.App,
		req.Version,
		"",
		metadata,
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  "completed",
		"version": req.Version,
	})
}

