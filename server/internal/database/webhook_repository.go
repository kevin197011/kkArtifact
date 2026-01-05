// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package database

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
)

// WebhookRepository handles webhook database operations
type WebhookRepository struct {
	db *DB
}

// NewWebhookRepository creates a new webhook repository
func NewWebhookRepository(db *DB) *WebhookRepository {
	return &WebhookRepository{db: db}
}

// Create creates a new webhook
func (r *WebhookRepository) Create(name string, eventTypes []string, url string, headers map[string]string, enabled bool, projectID, appID *int) (*Webhook, error) {
	var webhook Webhook
	var headersJSON sql.NullString
	
	if headers != nil {
		headersBytes, err := json.Marshal(headers)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal headers: %w", err)
		}
		headersJSON = sql.NullString{String: string(headersBytes), Valid: true}
	}

	query := `INSERT INTO webhooks (name, event_types, url, headers, enabled, project_id, app_id)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)
	          RETURNING id, name, event_types, url, headers, enabled, project_id, app_id, created_at`
	
	err := r.db.QueryRow(
		query,
		name,
		pq.Array(eventTypes),
		url,
		headersJSON,
		enabled,
		toNullInt64(projectID),
		toNullInt64(appID),
	).Scan(
		&webhook.ID,
		&webhook.Name,
		pq.Array(&webhook.EventTypes),
		&webhook.URL,
		&webhook.Headers,
		&webhook.Enabled,
		&webhook.ProjectID,
		&webhook.AppID,
		&webhook.CreatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}
	return &webhook, nil
}

// List lists all webhooks
func (r *WebhookRepository) List() ([]*Webhook, error) {
	query := `SELECT id, name, event_types, url, headers, enabled, project_id, app_id, created_at
	          FROM webhooks WHERE enabled = true ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list webhooks: %w", err)
	}
	defer rows.Close()

	var webhooks []*Webhook
	for rows.Next() {
		var webhook Webhook
		if err := rows.Scan(
			&webhook.ID,
			&webhook.Name,
			pq.Array(&webhook.EventTypes),
			&webhook.URL,
			&webhook.Headers,
			&webhook.Enabled,
			&webhook.ProjectID,
			&webhook.AppID,
			&webhook.CreatedAt,
		); err != nil {
			return nil, err
		}
		webhooks = append(webhooks, &webhook)
	}
	return webhooks, rows.Err()
}

// GetByID gets a webhook by ID
func (r *WebhookRepository) GetByID(id int) (*Webhook, error) {
	var webhook Webhook
	query := `SELECT id, name, event_types, url, headers, enabled, project_id, app_id, created_at
	          FROM webhooks WHERE id = $1`
	
	err := r.db.QueryRow(query, id).Scan(
		&webhook.ID,
		&webhook.Name,
		pq.Array(&webhook.EventTypes),
		&webhook.URL,
		&webhook.Headers,
		&webhook.Enabled,
		&webhook.ProjectID,
		&webhook.AppID,
		&webhook.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook: %w", err)
	}
	return &webhook, nil
}

// Update updates a webhook
func (r *WebhookRepository) Update(id int, name string, eventTypes []string, url string, headers map[string]string, enabled bool, projectID, appID *int) error {
	var headersJSON sql.NullString
	if headers != nil {
		headersBytes, err := json.Marshal(headers)
		if err != nil {
			return fmt.Errorf("failed to marshal headers: %w", err)
		}
		headersJSON = sql.NullString{String: string(headersBytes), Valid: true}
	}

	query := `UPDATE webhooks SET name = $1, event_types = $2, url = $3, headers = $4, enabled = $5, project_id = $6, app_id = $7
	          WHERE id = $8`
	_, err := r.db.Exec(query, name, pq.Array(eventTypes), url, headersJSON, enabled, toNullInt64(projectID), toNullInt64(appID), id)
	return err
}

// Delete deletes a webhook
func (r *WebhookRepository) Delete(id int) error {
	query := `DELETE FROM webhooks WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// FindByEventType finds webhooks that match the event type and optionally project/app
func (r *WebhookRepository) FindByEventType(eventType string, projectID, appID *int) ([]*Webhook, error) {
	query := `SELECT id, name, event_types, url, headers, enabled, project_id, app_id, created_at
	          FROM webhooks 
	          WHERE enabled = true 
	          AND $1 = ANY(event_types)
	          AND (project_id IS NULL OR project_id = $2)
	          AND (app_id IS NULL OR (project_id = $2 AND app_id = $3))
	          ORDER BY created_at DESC`
	
	var projectIDVal, appIDVal interface{}
	if projectID != nil {
		projectIDVal = *projectID
	} else {
		projectIDVal = nil
	}
	if appID != nil {
		appIDVal = *appID
	} else {
		appIDVal = nil
	}

	rows, err := r.db.Query(query, eventType, projectIDVal, appIDVal)
	if err != nil {
		return nil, fmt.Errorf("failed to find webhooks: %w", err)
	}
	defer rows.Close()

	var webhooks []*Webhook
	for rows.Next() {
		var webhook Webhook
		if err := rows.Scan(
			&webhook.ID,
			&webhook.Name,
			pq.Array(&webhook.EventTypes),
			&webhook.URL,
			&webhook.Headers,
			&webhook.Enabled,
			&webhook.ProjectID,
			&webhook.AppID,
			&webhook.CreatedAt,
		); err != nil {
			return nil, err
		}
		webhooks = append(webhooks, &webhook)
	}
	return webhooks, rows.Err()
}

