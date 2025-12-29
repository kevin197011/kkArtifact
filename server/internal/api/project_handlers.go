// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ProjectResponse represents a project in API response
type ProjectResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"` // RFC3339 format
}

// AppResponse represents an app in API response
type AppResponse struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"project_id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"` // RFC3339 format
}

// VersionResponse represents a version in API response
type VersionResponse struct {
	ID        int    `json:"id"`
	AppID     int    `json:"app_id"`
	Version   string `json:"version"`    // Version identifier (same as hash in database)
	CreatedAt string `json:"created_at"` // RFC3339 format
}

// handleListProjects godoc
// @Summary      List projects
// @Description  Get a list of all projects with pagination
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        limit   query     int  false  "Limit number of results (default: 50)"
// @Param        offset  query     int  false  "Offset for pagination (default: 0)"
// @Success      200     {array}   ProjectResponse
// @Failure      401     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Security     Bearer
// @Router       /projects [get]
func (h *Handler) handleListProjects(c *gin.Context) {
	limit := getIntQuery(c, "limit", 50)
	offset := getIntQuery(c, "offset", 0)

	projects, err := h.projectRepo.List(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format with properly formatted dates
	responses := make([]ProjectResponse, len(projects))
	for i, p := range projects {
		responses[i] = ProjectResponse{
			ID:        p.ID,
			Name:      p.Name,
			CreatedAt: p.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, responses)
}

func (h *Handler) handleListApps(c *gin.Context) {
	projectName := c.Param("project")
	limit := getIntQuery(c, "limit", 50)
	offset := getIntQuery(c, "offset", 0)

	project, err := h.projectRepo.CreateOrGet(projectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	apps, err := h.appRepo.ListByProject(project.ID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format with properly formatted dates
	responses := make([]AppResponse, len(apps))
	for i, a := range apps {
		responses[i] = AppResponse{
			ID:        a.ID,
			ProjectID: a.ProjectID,
			Name:      a.Name,
			CreatedAt: a.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, responses)
}

// handleListVersions godoc
// @Summary      List versions
// @Description  Get a list of all versions for a specific app
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        project  path      string  true   "Project name"
// @Param        app      path      string  true   "App name"
// @Param        limit    query     int     false  "Limit number of results (default: 50)"
// @Param        offset   query     int     false  "Offset for pagination (default: 0)"
// @Success      200      {array}   VersionResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Security     Bearer
// @Router       /projects/{project}/apps/{app}/versions [get]
func (h *Handler) handleListVersions(c *gin.Context) {
	projectName := c.Param("project")
	appName := c.Param("app")
	limit := getIntQuery(c, "limit", 50)
	offset := getIntQuery(c, "offset", 0)

	project, err := h.projectRepo.CreateOrGet(projectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	app, err := h.appRepo.CreateOrGet(project.ID, appName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	versions, err := h.versionRepo.ListByApp(app.ID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format with properly formatted dates
	responses := make([]VersionResponse, len(versions))
	for i, v := range versions {
		responses[i] = VersionResponse{
			ID:        v.ID,
			AppID:     v.AppID,
			Version:   v.Hash, // Map hash field from database to version in API
			CreatedAt: v.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, responses)
}

// handlePublicListProjects is a public (unauthenticated) version of handleListProjects
// It provides the same functionality but doesn't require authentication
func (h *Handler) handlePublicListProjects(c *gin.Context) {
	h.handleListProjects(c)
}

// handlePublicListApps is a public (unauthenticated) version of handleListApps
// It provides the same functionality but doesn't require authentication
func (h *Handler) handlePublicListApps(c *gin.Context) {
	h.handleListApps(c)
}

// handlePublicListVersions is a public (unauthenticated) version of handleListVersions
// It provides the same functionality but doesn't require authentication
func (h *Handler) handlePublicListVersions(c *gin.Context) {
	h.handleListVersions(c)
}
