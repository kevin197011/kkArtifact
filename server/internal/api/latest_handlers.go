// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kk/kkartifact-server/internal/database"
)

// handleGetLatestVersion returns the latest published version for an app
// handleGetLatestVersion godoc
// @Summary      Get latest published version
// @Description  Get the latest published version hash for a project/app
// @Tags         artifacts
// @Produce      json
// @Param        project  path  string  true  "Project name"
// @Param        app      path  string  true  "App name"
// @Success      200      {object}  map[string]string
// @Failure      404      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Security     Bearer
// @Router       /projects/{project}/apps/{app}/latest [get]
func (h *Handler) handleGetLatestVersion(c *gin.Context) {
	projectName := c.Param("project")
	appName := c.Param("app")

	// Get project and app
	project, err := h.projectRepo.GetByName(projectName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	app, err := h.appRepo.GetByName(project.ID, appName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "app not found"})
		return
	}

	// Get latest published version
	versionRepo := database.NewVersionRepository(h.db)
	latestVersion, err := versionRepo.GetLatestPublished(app.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no published version found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"project": projectName,
		"app":     appName,
		"version": latestVersion.Hash,
	})
}

