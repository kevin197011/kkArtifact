// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"time"

	"github.com/kk/kkartifact-server/internal/database"
)

func (h *Handler) webhookToResponse(webhook *database.Webhook) WebhookResponse {
	var headers *string
	if webhook.Headers.Valid {
		headers = &webhook.Headers.String
	}

	var projectID, appID *int
	var projectName, appName *string

	if webhook.ProjectID.Valid {
		pid := int(webhook.ProjectID.Int64)
		projectID = &pid
		projectName = h.lookupProjectName(pid)
	}

	if webhook.AppID.Valid {
		aid := int(webhook.AppID.Int64)
		appID = &aid
		appName = h.lookupAppName(aid)
	}

	return WebhookResponse{
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

func (h *Handler) lookupProjectName(id int) *string {
	var name string
	if err := h.db.QueryRow(`SELECT name FROM projects WHERE id = $1`, id).Scan(&name); err != nil {
		return nil
	}
	return &name
}

func (h *Handler) lookupAppName(id int) *string {
	var name string
	if err := h.db.QueryRow(`SELECT name FROM apps WHERE id = $1`, id).Scan(&name); err != nil {
		return nil
	}
	return &name
}
