// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kk/kkartifact-server/internal/auth"
	"github.com/kk/kkartifact-server/internal/database"
)

// PublishRequest represents a publish request
type PublishRequest struct {
	Project string `json:"project" binding:"required"`
	App     string `json:"app" binding:"required"`
	Version string `json:"version" binding:"required"`
}

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

	app, ok := h.resolvePublishApp(c, req)
	if !ok {
		return
	}

	if _, err := h.artifactManager.GetManifest(c.Request.Context(), req.Project, req.App, req.Version); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "version not found in storage"})
		return
	}

	if err := h.versionRepo.UnpublishAllVersions(app.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.versionRepo.SetPublished(app.ID, req.Version, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.publishEventWithContext(c, "publish", req.Project, req.App, req.Version, "", map[string]interface{}{
		"target_version": req.Version,
	})

	c.JSON(http.StatusOK, gin.H{
		"status":  "published",
		"project": req.Project,
		"app":     req.App,
		"version": req.Version,
	})
}

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

	app, ok := h.resolvePublishApp(c, req)
	if !ok {
		return
	}

	if err := h.versionRepo.SetPublished(app.ID, req.Version, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.publishEventWithContext(c, "unpublish", req.Project, req.App, req.Version, "", map[string]interface{}{
		"target_version": req.Version,
	})

	c.JSON(http.StatusOK, gin.H{
		"status":  "unpublished",
		"project": req.Project,
		"app":     req.App,
		"version": req.Version,
	})
}

// resolvePublishApp validates promote permission and resolves app + version for publish/unpublish.
func (h *Handler) resolvePublishApp(c *gin.Context, req PublishRequest) (*database.App, bool) {
	if !h.requireArtifactAccess(c, req.Project, req.App, auth.PermissionPromote) {
		return nil, false
	}

	project, err := h.projectRepo.GetByName(req.Project)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return nil, false
	}

	app, err := h.appRepo.GetByName(project.ID, req.App)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "app not found"})
		return nil, false
	}

	if _, err := h.versionRepo.GetByHash(app.ID, req.Version); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "version not found"})
		return nil, false
	}

	return app, true
}
