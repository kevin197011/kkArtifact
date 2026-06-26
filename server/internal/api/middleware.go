// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kk/kkartifact-server/internal/auth"
)

// adminMiddleware requires admin permission for all routes in the group.
func (h *Handler) adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !requireAdmin(c) {
			return
		}
		c.Next()
	}
}

// artifactPermissionMiddleware checks pull/push permission and token scope from path params.
func (h *Handler) artifactPermissionMiddleware(perm auth.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !h.requireArtifactAccess(c, c.Param("project"), c.Param("app"), perm) {
			return
		}
		c.Next()
	}
}

// deprecatedPublicListMiddleware marks legacy paginated public list endpoints as deprecated.
func deprecatedPublicListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Deprecation", "true")
		c.Header("Link", `</api/v1/public/inventory>; rel="successor-version"`)
		c.Next()
	}
}
