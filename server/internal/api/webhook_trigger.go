// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"encoding/json"
	"log"

	"github.com/kk/kkartifact-server/internal/database"
	"github.com/kk/kkartifact-server/internal/events"
)

func (h *Handler) triggerWebhook(webhook *database.Webhook, event *events.Event) {
	// Parse headers if present
	var headers map[string]string
	if webhook.Headers.Valid && webhook.Headers.String != "" {
		if err := json.Unmarshal([]byte(webhook.Headers.String), &headers); err != nil {
			log.Printf("Failed to parse webhook headers: %v", err)
			headers = make(map[string]string)
		}
	} else {
		headers = make(map[string]string)
	}

	// Convert database.Webhook to events-compatible format if needed
	// The event is already in the right format

	// Create webhook sender
	sender := events.NewWebhookSender()

	// Send webhook
	if err := sender.Send(webhook.URL, headers, event); err != nil {
		log.Printf("Failed to send webhook %d to %s: %v", webhook.ID, webhook.URL, err)
		// TODO: Store webhook failure in audit log
	} else {
		log.Printf("Webhook %d triggered successfully", webhook.ID)
	}
}

