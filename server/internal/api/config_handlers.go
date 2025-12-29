// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kk/kkartifact-server/internal/database"
)

// handleGetConfig gets global configuration
// handleGetConfig godoc
// @Summary      Get config
// @Description  Get the global configuration (e.g., version retention limit)
// @Tags         config
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     Bearer
// @Router       /config [get]
func (h *Handler) handleGetConfig(c *gin.Context) {
	configRepo := database.NewConfigRepository(h.db)
	
	// Get version retention limit
	retentionLimit, err := configRepo.Get("version_retention_limit")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	limit, _ := strconv.Atoi(retentionLimit)
	
	c.JSON(http.StatusOK, gin.H{
		"version_retention_limit": limit,
	})
}

// handleUpdateConfig updates global configuration
// handleUpdateConfig godoc
// @Summary      Update config
// @Description  Update the global configuration
// @Tags         config
// @Accept       json
// @Produce      json
// @Param        request  body      map[string]interface{}  true  "Configuration updates"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Security     Bearer
// @Router       /config [put]
func (h *Handler) handleUpdateConfig(c *gin.Context) {
	var req struct {
		VersionRetentionLimit *int `json:"version_retention_limit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configRepo := database.NewConfigRepository(h.db)

	if req.VersionRetentionLimit != nil {
		if *req.VersionRetentionLimit < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "version_retention_limit must be at least 1"})
			return
		}
		if err := configRepo.Set("version_retention_limit", strconv.Itoa(*req.VersionRetentionLimit)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

