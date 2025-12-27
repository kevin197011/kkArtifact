// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kk/kkartifact-server/internal/auth"
	"github.com/kk/kkartifact-server/internal/database"
	"github.com/kk/kkartifact-server/internal/storage"
)

// Handler handles HTTP requests
type Handler struct {
	db              *database.DB
	storage         storage.Storage
	artifactManager *storage.ArtifactManager
	authenticator   *auth.TokenAuthenticator
	projectRepo     *database.ProjectRepository
	appRepo         *database.AppRepository
	versionRepo     *database.VersionRepository
}

// NewHandler creates a new API handler
func NewHandler(db *database.DB, storageBackend storage.Storage, authenticator *auth.TokenAuthenticator) *Handler {
	return &Handler{
		db:              db,
		storage:         storageBackend,
		artifactManager: storage.NewArtifactManager(storageBackend),
		authenticator:   authenticator,
		projectRepo:     database.NewProjectRepository(db),
		appRepo:         database.NewAppRepository(db),
		versionRepo:     database.NewVersionRepository(db),
	}
}

// RegisterRoutes registers all API routes
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	
	// Public routes
	api.GET("/health", h.handleHealth)
	api.HEAD("/health", h.handleHealth) // Support HEAD requests for health checks
	
	// Login endpoint (public)
	api.POST("/login", h.handleLogin)
	
	// Token creation endpoint (public for initial setup, consider protecting in production)
	// Register BEFORE protected routes to avoid conflicts
	api.POST("/tokens", h.handleCreateToken)
	
	// Protected routes
	protected := api.Group("")
	protected.Use(h.authenticator.AuthMiddleware())
	{
		// List endpoints
		protected.GET("/projects", h.handleListProjects)
		protected.GET("/projects/:project/apps", h.handleListApps)
		protected.GET("/projects/:project/apps/:app/versions", h.handleListVersions)
		protected.GET("/manifest/:project/:app/:hash", h.handleGetManifest)
		protected.GET("/file/:project/:app/:hash", h.handleGetFile)
		
		// Upload endpoints
		protected.POST("/upload/init", h.handleInitUpload)
		protected.POST("/file/:project/:app/:hash", h.handleUploadFile)
		protected.POST("/upload/finish", h.handleFinishUpload)
		
		// Webhook endpoints
		protected.POST("/webhooks", h.handleCreateWebhook)
		protected.GET("/webhooks", h.handleListWebhooks)
		protected.GET("/webhooks/:id", h.handleGetWebhook)
		protected.PUT("/webhooks/:id", h.handleUpdateWebhook)
		protected.DELETE("/webhooks/:id", h.handleDeleteWebhook)
		
		// Config endpoints
		protected.GET("/config", h.handleGetConfig)
		protected.PUT("/config", h.handleUpdateConfig)
		
		// Promote endpoint
		protected.POST("/promote", h.handlePromote)
		
		// Audit logs endpoint
		protected.GET("/audit-logs", h.handleListAuditLogs)
		
		// Token management endpoints (list and delete require auth)
		protected.GET("/tokens", h.handleListTokens)
		protected.DELETE("/tokens/:id", h.handleDeleteToken)
		
		// Storage sync endpoint (admin only - rebuilds database from storage)
		protected.POST("/sync-storage", h.handleSyncStorage)
	}
}

func (h *Handler) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// handleListProjects, handleListApps, handleListVersions are moved to project_handlers.go

// handleGetManifest is moved to manifest_handlers.go

func (h *Handler) handleGetFile(c *gin.Context) {
	project := c.Param("project")
	app := c.Param("app")
	hash := c.Param("hash")
	filePath := c.Query("path")
	
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path parameter is required"})
		return
	}
	
	fullPath := project + "/" + app + "/" + hash + "/" + filePath
	
	// Handle HEAD request for file existence check
	if c.Request.Method == "HEAD" {
		exists, err := h.storage.Exists(c.Request.Context(), fullPath)
		if err != nil || !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.Status(http.StatusOK)
		return
	}
	
	reader, err := h.storage.Get(c.Request.Context(), fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	defer reader.Close()
	
	// Support HTTP Range requests for resume download
	rangeHeader := c.GetHeader("Range")
	if rangeHeader != "" {
		// Parse range header (e.g., "bytes=0-1023" or "bytes=1024-")
		var start, end int64 = 0, -1
		if len(rangeHeader) > 6 && rangeHeader[:6] == "bytes=" {
			rangeStr := rangeHeader[6:]
			parts := strings.Split(rangeStr, "-")
			if len(parts) == 2 {
				fmt.Sscanf(parts[0], "%d", &start)
				if parts[1] != "" {
					fmt.Sscanf(parts[1], "%d", &end)
				}
			}
		}
		
		// Try to seek to start position if reader supports it
		if seeker, ok := reader.(io.Seeker); ok {
			if _, err := seeker.Seek(start, io.SeekStart); err == nil {
				// Set Content-Range header for partial content
				c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/*", start, end))
				c.Header("Accept-Ranges", "bytes")
				if end > 0 {
					length := end - start + 1
					c.DataFromReader(http.StatusPartialContent, length, "application/octet-stream", reader, nil)
				} else {
					c.DataFromReader(http.StatusPartialContent, -1, "application/octet-stream", reader, nil)
				}
				return
			}
		}
	}
	
	// Record audit log for file download (pull operation)
	auditRepo := database.NewAuditRepository(h.db)
	projectObj, err := h.projectRepo.CreateOrGet(project)
	if err == nil {
		var projectID *int
		var appID *int
		if projectObj != nil {
			projectID = &projectObj.ID
			appObj, err := h.appRepo.CreateOrGet(projectObj.ID, app)
			if err == nil && appObj != nil {
				appID = &appObj.ID
			}
		}
		metadata := map[string]interface{}{
			"file_path": filePath,
		}
		_ = auditRepo.Create("pull", projectID, appID, hash, "", metadata)
	}

	// Full file download
	c.Header("Accept-Ranges", "bytes")
	c.DataFromReader(http.StatusOK, -1, "application/octet-stream", reader, nil)
}

// authMiddleware is moved to auth.AuthMiddleware() which supports both JWT and API tokens

// getIntQuery is moved to helpers.go

