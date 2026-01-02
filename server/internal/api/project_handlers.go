// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kk/kkartifact-server/internal/database"
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

// handleDeleteProject godoc
// @Summary      Delete project
// @Description  Delete a project and all its apps and versions (cascade delete)
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        project  path      string  true   "Project name"
// @Success      200      {object}  map[string]string
// @Failure      401      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Security     Bearer
// @Router       /projects/{project} [delete]
func (h *Handler) handleDeleteProject(c *gin.Context) {
	projectName := c.Param("project")

	// Get project to verify it exists
	project, err := h.projectRepo.CreateOrGet(projectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete project from database (cascade will delete apps and versions)
	if err := h.projectRepo.Delete(projectName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete project: %v", err)})
		return
	}

	// Delete project from storage
	if err := h.artifactManager.DeleteProject(c.Request.Context(), projectName); err != nil {
		// Log error but don't fail - database is already deleted
		// This is a best-effort cleanup
		_ = err
	}

	// Create audit log
	auditRepo := database.NewAuditRepository(h.db)
	projectID := project.ID
	_ = auditRepo.Create("project_delete", &projectID, nil, "", "", map[string]interface{}{
		"project_name": projectName,
	})

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// handleDeleteApp godoc
// @Summary      Delete app
// @Description  Delete an app and all its versions (cascade delete)
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        project  path      string  true   "Project name"
// @Param        app      path      string  true   "App name"
// @Success      200      {object}  map[string]string
// @Failure      401      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Security     Bearer
// @Router       /projects/{project}/apps/{app} [delete]
func (h *Handler) handleDeleteApp(c *gin.Context) {
	projectName := c.Param("project")
	appName := c.Param("app")

	// Get project and app to verify they exist
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

	// Delete app from database (cascade will delete versions)
	if err := h.appRepo.Delete(project.ID, appName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete app: %v", err)})
		return
	}

	// Delete app from storage
	if err := h.artifactManager.DeleteApp(c.Request.Context(), projectName, appName); err != nil {
		// Log error but don't fail - database is already deleted
		_ = err
	}

	// Create audit log
	auditRepo := database.NewAuditRepository(h.db)
	projectID := project.ID
	appID := app.ID
	_ = auditRepo.Create("app_delete", &projectID, &appID, "", "", map[string]interface{}{
		"project_name": projectName,
		"app_name":      appName,
	})

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// handleDeleteVersion godoc
// @Summary      Delete version
// @Description  Delete a version
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        project  path      string  true   "Project name"
// @Param        app      path      string  true   "App name"
// @Param        version  path      string  true   "Version hash"
// @Success      200      {object}  map[string]string
// @Failure      401      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Security     Bearer
// @Router       /projects/{project}/apps/{app}/versions/{version} [delete]
func (h *Handler) handleDeleteVersion(c *gin.Context) {
	projectName := c.Param("project")
	appName := c.Param("app")
	versionHash := c.Param("version")

	// Get project and app to verify they exist
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

	// Delete version from database
	if err := h.versionRepo.Delete(app.ID, versionHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete version: %v", err)})
		return
	}

	// Delete version from storage
	if err := h.artifactManager.DeleteVersion(c.Request.Context(), projectName, appName, versionHash); err != nil {
		// Log error but don't fail - database is already deleted
		_ = err
	}

	// Create audit log
	auditRepo := database.NewAuditRepository(h.db)
	projectID := project.ID
	appID := app.ID
	_ = auditRepo.Create("version_delete", &projectID, &appID, versionHash, "", map[string]interface{}{
		"project_name": projectName,
		"app_name":      appName,
		"version":       versionHash,
	})

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
