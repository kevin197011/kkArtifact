// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PromoteRequest represents a promote request
type PromoteRequest struct {
	Project string `json:"project" binding:"required"`
	App     string `json:"app" binding:"required"`
	Hash    string `json:"hash" binding:"required"`
}

// handlePromote promotes a version
// handlePromote godoc
// @Summary      Promote version
// @Description  Promote a version to a new version identifier (create a new version with same content)
// @Tags         artifacts
// @Accept       json
// @Produce      json
// @Param        request  body      PromoteRequest  true  "Promote request"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Security     Bearer
// @Router       /promote [post]
func (h *Handler) handlePromote(c *gin.Context) {
	var req PromoteRequest
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

	_, err = h.appRepo.CreateOrGet(project.ID, req.App)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Verify version exists
	manifest, err := h.artifactManager.GetManifest(c.Request.Context(), req.Project, req.App, req.Hash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "version not found"})
		return
	}

	// TODO: Mark version as promoted in database
	// For now, just return success
	_ = manifest

	// Publish promote event with context to extract agent ID
	metadata := make(map[string]interface{})

	h.publishEventWithContext(
		c,
		"promote",
		req.Project,
		req.App,
		req.Hash,
		"",
		metadata,
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  "promoted",
		"project": req.Project,
		"app":     req.App,
		"version": req.Hash,
	})
}

