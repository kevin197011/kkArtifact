// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kk/kkartifact-server/internal/auth"
)

// Helper functions for API handlers

// getIntQuery gets an integer query parameter
func getIntQuery(c *gin.Context, key string, defaultValue int) int {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	if intValue < 1 {
		return defaultValue
	}
	if intValue > 500 {
		return 500
	}
	return intValue
}

// getIntParam gets an integer path parameter
func getIntParam(c *gin.Context, key string, defaultValue int) int {
	value := c.Param(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// isAdminRequest returns true if the caller is an admin (JWT session or API token with admin permission)
func isAdminRequest(c *gin.Context) bool {
	return callerHasPermission(c, auth.PermissionAdmin)
}

// requireAdmin aborts with 403 unless the caller is an admin
func requireAdmin(c *gin.Context) bool {
	if isAdminRequest(c) {
		return true
	}
	c.JSON(http.StatusForbidden, gin.H{"error": "admin permission required"})
	c.Abort()
	return false
}

// callerHasPermission returns true when the caller may perform the required action.
// JWT admin sessions have all permissions; other JWT sessions may only pull.
// API tokens are checked against their permission list.
func callerHasPermission(c *gin.Context, required auth.Permission) bool {
	if sessionInfo, exists := c.Get("session_info"); exists {
		if s, ok := sessionInfo.(*auth.SessionInfo); ok {
			if s.IsAdmin {
				return true
			}
			return required == auth.PermissionPull
		}
	}
	if tokenInfo, exists := c.Get("token_info"); exists {
		if t, ok := tokenInfo.(*auth.TokenInfo); ok {
			return auth.HasPermission(t.Permissions, required)
		}
	}
	return false
}

// requireArtifactAccess checks permission and API token project/app scope.
func (h *Handler) requireArtifactAccess(c *gin.Context, project, app string, perm auth.Permission) bool {
	if !callerHasPermission(c, perm) {
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("%s permission required", perm)})
		c.Abort()
		return false
	}
	return h.checkTokenScope(c, project, app)
}

// checkTokenScope enforces project/app scope for API tokens. JWT sessions are unrestricted.
func (h *Handler) checkTokenScope(c *gin.Context, projectName, appName string) bool {
	if _, exists := c.Get("session_info"); exists {
		return true
	}

	tokenInfo, exists := c.Get("token_info")
	if !exists {
		return true
	}

	t, ok := tokenInfo.(*auth.TokenInfo)
	if !ok || t.ProjectID == nil {
		return true
	}

	project, err := h.projectRepo.GetByName(projectName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		c.Abort()
		return false
	}
	if project.ID != *t.ProjectID {
		c.JSON(http.StatusForbidden, gin.H{"error": "token scope does not include this project"})
		c.Abort()
		return false
	}

	if t.AppID == nil {
		return true
	}

	app, err := h.appRepo.GetByName(project.ID, appName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "app not found"})
		c.Abort()
		return false
	}
	if app.ID != *t.AppID {
		c.JSON(http.StatusForbidden, gin.H{"error": "token scope does not include this app"})
		c.Abort()
		return false
	}

	return true
}
