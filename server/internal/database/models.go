// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package database

import (
	"database/sql"
	"time"
)

// Project represents a project in the database
type Project struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

// App represents an app in the database
type App struct {
	ID        int       `db:"id"`
	ProjectID int       `db:"project_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

// Version represents a version in the database
type Version struct {
	ID          int       `db:"id"`
	AppID       int       `db:"app_id"`
	Hash        string    `db:"hash"`
	IsPublished bool      `db:"is_published"`
	CreatedAt   time.Time `db:"created_at"`
}

// Token represents a token in the database
type Token struct {
	ID         int            `db:"id"`
	TokenHash  string         `db:"token_hash"`
	Name       sql.NullString `db:"name"`
	ProjectID  sql.NullInt64  `db:"project_id"`
	AppID      sql.NullInt64  `db:"app_id"`
	Permissions []string      `db:"permissions"`
	ExpiresAt  sql.NullTime   `db:"expires_at"`
	CreatedAt  time.Time      `db:"created_at"`
}

// Webhook represents a webhook in the database
type Webhook struct {
	ID         int            `db:"id"`
	Name       string         `db:"name"`
	EventTypes []string       `db:"event_types"`
	URL        string         `db:"url"`
	Headers    sql.NullString `db:"headers"`
	Enabled    bool           `db:"enabled"`
	ProjectID  sql.NullInt64  `db:"project_id"`
	AppID      sql.NullInt64  `db:"app_id"`
	CreatedAt  time.Time      `db:"created_at"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID          int            `db:"id"`
	Operation   string         `db:"operation"`
	ProjectID   sql.NullInt64  `db:"project_id"`
	AppID       sql.NullInt64  `db:"app_id"`
	VersionHash sql.NullString `db:"version_hash"`
	AgentID     sql.NullString `db:"agent_id"`
	Metadata    sql.NullString `db:"metadata"`
	CreatedAt   time.Time      `db:"created_at"`
}

// Config represents a global configuration entry
type Config struct {
	ID        int       `db:"id"`
	Key       string    `db:"key"`
	Value     string    `db:"value"`
	UpdatedAt time.Time `db:"updated_at"`
}

