// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"time"

	"github.com/kk/kkartifact-server/internal/database"
	"github.com/kk/kkartifact-server/internal/events"
)

// publishEvent publishes an event and triggers webhooks
func (h *Handler) publishEvent(eventType events.EventType, project, app, version, agentID string, metadata map[string]interface{}) {
	event := &events.Event{
		Type:      eventType,
		Project:   project,
		App:       app,
		Version:   version,
		AgentID:   agentID,
		Metadata:  metadata,
		Timestamp: time.Now(),
	}

	// TODO: Integrate with event bus
	// For now, just log the event
	_ = event

	// Get webhooks and trigger them
	webhookRepo := database.NewWebhookRepository(h.db)
	
	// Get project and app IDs
	projectObj, err := h.projectRepo.CreateOrGet(project)
	if err != nil {
		return
	}
	
	var projectID *int
	var appID *int
	if projectObj != nil {
		projectID = &projectObj.ID
		
		appObj, err := h.appRepo.CreateOrGet(projectObj.ID, app)
		if err == nil && appObj != nil {
			appID = &appObj.ID
		}
	}

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

