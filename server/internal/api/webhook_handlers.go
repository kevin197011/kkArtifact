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

// CreateWebhookRequest represents a webhook creation request
type CreateWebhookRequest struct {
	Name       string            `json:"name" binding:"required"`
	EventTypes []string          `json:"event_types" binding:"required"`
	URL        string            `json:"url" binding:"required,url"`
	Headers    map[string]string `json:"headers,omitempty"`
	Enabled    bool              `json:"enabled"`
	ProjectID  *int              `json:"project_id,omitempty"`
	AppID      *int              `json:"app_id,omitempty"`
}

// handleCreateWebhook creates a new webhook
func (h *Handler) handleCreateWebhook(c *gin.Context) {
	var req CreateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	webhookRepo := database.NewWebhookRepository(h.db)
	webhook, err := webhookRepo.Create(
		req.Name,
		req.EventTypes,
		req.URL,
		req.Headers,
		req.Enabled,
		req.ProjectID,
		req.AppID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, webhook)
}

// handleListWebhooks lists all webhooks
func (h *Handler) handleListWebhooks(c *gin.Context) {
	webhookRepo := database.NewWebhookRepository(h.db)
	webhooks, err := webhookRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, webhooks)
}

// handleGetWebhook gets a webhook by ID
func (h *Handler) handleGetWebhook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid webhook ID"})
		return
	}

	webhookRepo := database.NewWebhookRepository(h.db)
	webhook, err := webhookRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if webhook == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "webhook not found"})
		return
	}

	c.JSON(http.StatusOK, webhook)
}

// handleUpdateWebhook updates a webhook
func (h *Handler) handleUpdateWebhook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid webhook ID"})
		return
	}

	var req struct {
		Name       string            `json:"name"`
		EventTypes []string          `json:"event_types"`
		URL        string            `json:"url"`
		Headers    map[string]string `json:"headers"`
		Enabled    *bool             `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	webhookRepo := database.NewWebhookRepository(h.db)
	webhook, err := webhookRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if webhook == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "webhook not found"})
		return
	}

	// Update fields if provided
	name := req.Name
	if name == "" {
		name = webhook.Name
	}
	eventTypes := req.EventTypes
	if len(eventTypes) == 0 {
		eventTypes = webhook.EventTypes
	}
	url := req.URL
	if url == "" {
		url = webhook.URL
	}
	enabled := webhook.Enabled
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	if err := webhookRepo.Update(id, name, eventTypes, url, req.Headers, enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updated, err := webhookRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// handleDeleteWebhook deletes a webhook
func (h *Handler) handleDeleteWebhook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid webhook ID"})
		return
	}

	webhookRepo := database.NewWebhookRepository(h.db)
	if err := webhookRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

