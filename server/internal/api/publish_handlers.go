// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PublishRequest represents a publish request
type PublishRequest struct {
	Project string `json:"project" binding:"required"`
	App     string `json:"app" binding:"required"`
	Version string `json:"version" binding:"required"`
}

// handlePublish publishes a version (marks it as published)
// handlePublish godoc
// @Summary      Publish version
// @Description  Mark a version as published so it can be retrieved via pull latest
// @Tags         artifacts
// @Accept       json
// @Produce      json
// @Param        request  body      PublishRequest  true  "Publish request"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Security     Bearer
// @Router       /publish [post]
func (h *Handler) handlePublish(c *gin.Context) {
	var req PublishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get project and app
	project, err := h.projectRepo.GetByName(req.Project)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	app, err := h.appRepo.GetByName(project.ID, req.App)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "app not found"})
		return
	}

	// Verify version exists in database
	versions, err := h.versionRepo.ListByApp(app.ID, 10000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var versionExists bool
	for _, v := range versions {
		if v.Hash == req.Version {
			versionExists = true
			break
		}
	}

	if !versionExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "version not found"})
		return
	}

	// Verify version exists in storage
	_, err = h.artifactManager.GetManifest(c.Request.Context(), req.Project, req.App, req.Version)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "version not found in storage"})
		return
	}

	// Unpublish all other versions for this app (only one version can be published at a time)
	if err := h.versionRepo.UnpublishAllVersions(app.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Mark version as published
	if err := h.versionRepo.SetPublished(app.ID, req.Version, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Publish publish event with context to extract agent ID
	metadata := make(map[string]interface{})
	metadata["target_version"] = req.Version

	h.publishEventWithContext(
		c,
		"publish",
		req.Project,
		req.App,
		req.Version,
		"",
		metadata,
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  "published",
		"project": req.Project,
		"app":     req.App,
		"version": req.Version,
	})
}

// handleUnpublish unpublishes a version (marks it as not published)
// handleUnpublish godoc
// @Summary      Unpublish version
// @Description  Mark a version as unpublished
// @Tags         artifacts
// @Accept       json
// @Produce      json
// @Param        request  body      PublishRequest  true  "Unpublish request"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Security     Bearer
// @Router       /unpublish [post]
func (h *Handler) handleUnpublish(c *gin.Context) {
	var req PublishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get project and app
	project, err := h.projectRepo.GetByName(req.Project)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	app, err := h.appRepo.GetByName(project.ID, req.App)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "app not found"})
		return
	}

	// Verify version exists in database
	versions, err := h.versionRepo.ListByApp(app.ID, 10000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var versionExists bool
	for _, v := range versions {
		if v.Hash == req.Version {
			versionExists = true
			break
		}
	}

	if !versionExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "version not found"})
		return
	}

	// Mark version as unpublished
	if err := h.versionRepo.SetPublished(app.ID, req.Version, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Publish unpublish event with context to extract agent ID
	metadata := make(map[string]interface{})
	metadata["target_version"] = req.Version

	h.publishEventWithContext(
		c,
		"unpublish",
		req.Project,
		req.App,
		req.Version,
		"",
		metadata,
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  "unpublished",
		"project": req.Project,
		"app":     req.App,
		"version": req.Version,
	})
}
