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

// handleCreateWebhook godoc
// @Summary      Create webhook
// @Description  Create a new webhook for event notifications (can be global, project-level, or app-level)
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        request  body      CreateWebhookRequest  true  "Webhook creation request"
// @Success      201      {object}  WebhookResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      403      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Security     Bearer
// @Router       /webhooks [post]
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

	// Convert to response format with project and app names
	c.JSON(http.StatusCreated, h.webhookToResponse(webhook))
}

// WebhookResponse represents a webhook in API response
type WebhookResponse struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	EventTypes  []string `json:"event_types"`
	URL         string   `json:"url"`
	Headers     *string  `json:"headers,omitempty"`
	Enabled     bool     `json:"enabled"`
	ProjectID   *int     `json:"project_id,omitempty"`
	AppID       *int     `json:"app_id,omitempty"`
	ProjectName *string  `json:"project_name,omitempty"`
	AppName     *string  `json:"app_name,omitempty"`
	CreatedAt   string   `json:"created_at"`
}

// handleListWebhooks godoc
// @Summary      List webhooks
// @Description  Get a list of all webhooks (including disabled ones)
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Success      200  {array}   WebhookResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     Bearer
// @Router       /webhooks [get]
func (h *Handler) handleListWebhooks(c *gin.Context) {
	webhookRepo := database.NewWebhookRepository(h.db)
	webhooks, err := webhookRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format with project and app names
	responses := make([]WebhookResponse, len(webhooks))
	for i, webhook := range webhooks {
		responses[i] = h.webhookToResponse(webhook)
	}

	c.JSON(http.StatusOK, responses)
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

	// Convert to response format with project and app names
	c.JSON(http.StatusOK, h.webhookToResponse(webhook))
}

// handleUpdateWebhook godoc
// @Summary      Update webhook
// @Description  Update an existing webhook
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        id       path      int                    true  "Webhook ID"
// @Param        request  body      CreateWebhookRequest   true  "Webhook update request"
// @Success      200      {object}  WebhookResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      403      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Security     Bearer
// @Router       /webhooks/{id} [put]
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
		ProjectID  *int              `json:"project_id,omitempty"`
		AppID      *int              `json:"app_id,omitempty"`
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

	var projectID, appID *int
	if req.ProjectID != nil {
		projectID = req.ProjectID
	} else {
		// Keep existing values if not provided
		if webhook.ProjectID.Valid {
			pid := int(webhook.ProjectID.Int64)
			projectID = &pid
		}
	}
	if req.AppID != nil {
		appID = req.AppID
	} else {
		// Keep existing values if not provided
		if webhook.AppID.Valid {
			aid := int(webhook.AppID.Int64)
			appID = &aid
		}
	}

	if err := webhookRepo.Update(id, name, eventTypes, url, req.Headers, enabled, projectID, appID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updated, err := webhookRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.webhookToResponse(updated))
}

// handleDeleteWebhook godoc
// @Summary      Delete webhook
// @Description  Delete a webhook by ID
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Webhook ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     Bearer
// @Router       /webhooks/{id} [delete]
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
