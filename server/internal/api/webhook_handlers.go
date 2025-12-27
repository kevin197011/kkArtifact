// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"net/http"
	"strconv"
	"time"

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

	// Convert to response format with project and app names
	var projectID, appID *int
	var projectName, appName *string

	if webhook.ProjectID.Valid {
		pid := int(webhook.ProjectID.Int64)
		projectID = &pid
		// Get project name
		var name string
		query := `SELECT name FROM projects WHERE id = $1`
		if err := h.db.QueryRow(query, pid).Scan(&name); err == nil {
			projectName = &name
		}
	}

	if webhook.AppID.Valid {
		aid := int(webhook.AppID.Int64)
		appID = &aid
		// Get app name
		var name string
		query := `SELECT name FROM apps WHERE id = $1`
		if err := h.db.QueryRow(query, aid).Scan(&name); err == nil {
			appName = &name
		}
	}

	var headers *string
	if webhook.Headers.Valid {
		headers = &webhook.Headers.String
	}

	response := WebhookResponse{
		ID:          webhook.ID,
		Name:        webhook.Name,
		EventTypes:  webhook.EventTypes,
		URL:         webhook.URL,
		Headers:     headers,
		Enabled:     webhook.Enabled,
		ProjectID:   projectID,
		AppID:       appID,
		ProjectName: projectName,
		AppName:     appName,
		CreatedAt:   webhook.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, response)
}

// WebhookResponse represents a webhook in API response
type WebhookResponse struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	EventTypes []string `json:"event_types"`
	URL        string  `json:"url"`
	Headers    *string `json:"headers,omitempty"`
	Enabled    bool    `json:"enabled"`
	ProjectID  *int    `json:"project_id,omitempty"`
	AppID      *int    `json:"app_id,omitempty"`
	ProjectName *string `json:"project_name,omitempty"`
	AppName     *string `json:"app_name,omitempty"`
	CreatedAt  string  `json:"created_at"`
}

// handleListWebhooks lists all webhooks
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
		var projectID, appID *int
		var projectName, appName *string

		if webhook.ProjectID.Valid {
			pid := int(webhook.ProjectID.Int64)
			projectID = &pid
			// Get project name
			var name string
			query := `SELECT name FROM projects WHERE id = $1`
			if err := h.db.QueryRow(query, pid).Scan(&name); err == nil {
				projectName = &name
			}
		}

		if webhook.AppID.Valid {
			aid := int(webhook.AppID.Int64)
			appID = &aid
			// Get app name
			var name string
			query := `SELECT name FROM apps WHERE id = $1`
			if err := h.db.QueryRow(query, aid).Scan(&name); err == nil {
				appName = &name
			}
		}

		var headers *string
		if webhook.Headers.Valid {
			headers = &webhook.Headers.String
		}

		responses[i] = WebhookResponse{
			ID:          webhook.ID,
			Name:        webhook.Name,
			EventTypes:  webhook.EventTypes,
			URL:         webhook.URL,
			Headers:     headers,
			Enabled:     webhook.Enabled,
			ProjectID:   projectID,
			AppID:       appID,
			ProjectName: projectName,
			AppName:     appName,
			CreatedAt:   webhook.CreatedAt.Format(time.RFC3339),
		}
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
	var projectID, appID *int
	var projectName, appName *string

	if webhook.ProjectID.Valid {
		pid := int(webhook.ProjectID.Int64)
		projectID = &pid
		// Get project name
		var name string
		query := `SELECT name FROM projects WHERE id = $1`
		if err := h.db.QueryRow(query, pid).Scan(&name); err == nil {
			projectName = &name
		}
	}

	if webhook.AppID.Valid {
		aid := int(webhook.AppID.Int64)
		appID = &aid
		// Get app name
		var name string
		query := `SELECT name FROM apps WHERE id = $1`
		if err := h.db.QueryRow(query, aid).Scan(&name); err == nil {
			appName = &name
		}
	}

	var headers *string
	if webhook.Headers.Valid {
		headers = &webhook.Headers.String
	}

	response := WebhookResponse{
		ID:          webhook.ID,
		Name:        webhook.Name,
		EventTypes:  webhook.EventTypes,
		URL:         webhook.URL,
		Headers:     headers,
		Enabled:     webhook.Enabled,
		ProjectID:   projectID,
		AppID:       appID,
		ProjectName: projectName,
		AppName:     appName,
		CreatedAt:   webhook.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
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

