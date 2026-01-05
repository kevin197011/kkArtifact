// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kk/kkartifact-server/internal/database"
	"github.com/kk/kkartifact-server/internal/events"
)

// publishEvent publishes an event and triggers webhooks
func (h *Handler) publishEvent(eventType events.EventType, project, app, version, agentID string, metadata map[string]interface{}) {
	h.publishEventWithContext(nil, eventType, project, app, version, agentID, metadata)
}

// publishEventWithContext publishes an event with context for extracting client information
func (h *Handler) publishEventWithContext(c *gin.Context, eventType events.EventType, project, app, version, agentID string, metadata map[string]interface{}) {
	// Extract agent ID from context if not provided
	if agentID == "" && c != nil {
		agentID = getAgentIDFromRequest(c)
	}

	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	event := &events.Event{
		Type:      eventType,
		Project:   project,
		App:       app,
		Version:   version,
		AgentID:   agentID,
		Metadata:  metadata,
		Timestamp: time.Now(),
	}

	// Publish event to event bus if available
	if h.eventBus != nil {
		if err := h.eventBus.Publish(event); err != nil {
			// Log error but continue with audit log and webhook triggering
			log.Printf("Failed to publish event to event bus: %v", err)
		}
	}

	// Get project and app IDs
	projectObj, err := h.projectRepo.CreateOrGet(project)
	if err != nil {
		return
	}
	
	var projectID *int
	var appID *int
	if projectObj != nil {
		projectID = &projectObj.ID
		
		if app != "" {
			appObj, err := h.appRepo.CreateOrGet(projectObj.ID, app)
			if err == nil && appObj != nil {
				appID = &appObj.ID
			}
		}
	}

	// Record audit log
	auditRepo := database.NewAuditRepository(h.db)
	_ = auditRepo.Create(string(eventType), projectID, appID, version, agentID, metadata)

	// Get webhooks and trigger them
	webhookRepo := database.NewWebhookRepository(h.db)
	webhooks, err := webhookRepo.FindByEventType(string(eventType), projectID, appID)
	if err != nil {
		return
	}

	// Trigger webhooks asynchronously
	for _, webhook := range webhooks {
		go h.triggerWebhook(webhook, event)
	}
}

// triggerWebhook is in webhook_trigger.go

