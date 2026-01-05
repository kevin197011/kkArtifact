// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// AuditRepository handles audit log database operations
type AuditRepository struct {
	db *DB
}

// NewAuditRepository creates a new audit repository
func NewAuditRepository(db *DB) *AuditRepository {
	return &AuditRepository{db: db}
}

// Create creates a new audit log entry
func (r *AuditRepository) Create(operation string, projectID, appID *int, versionHash, agentID string, metadata map[string]interface{}) error {
	var metadataJSON sql.NullString
	if metadata != nil {
		metadataBytes, err := json.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metadataJSON = sql.NullString{String: string(metadataBytes), Valid: true}
	}

	query := `INSERT INTO audit_logs (operation, project_id, app_id, version_hash, agent_id, metadata)
	          VALUES ($1, $2, $3, $4, $5, $6)`
	
	_, err := r.db.Exec(
		query,
		operation,
		toNullInt64(projectID),
		toNullInt64(appID),
		toNullString(versionHash),
		toNullString(agentID),
		metadataJSON,
	)
	return err
}

// List lists audit logs with optional filters
func (r *AuditRepository) List(projectID, appID *int, limit, offset int) ([]*AuditLog, error) {
	query := `SELECT id, operation, project_id, app_id, version_hash, agent_id, metadata, created_at
	          FROM audit_logs 
	          WHERE ($1::int IS NULL OR project_id = $1)
	          AND ($2::int IS NULL OR app_id = $2)
	          ORDER BY created_at DESC 
	          LIMIT $3 OFFSET $4`
	
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

	rows, err := r.db.Query(query, projectIDVal, appIDVal, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit logs: %w", err)
	}
	defer rows.Close()

	var logs []*AuditLog
	for rows.Next() {
		var log AuditLog
		if err := rows.Scan(
			&log.ID,
			&log.Operation,
			&log.ProjectID,
			&log.AppID,
			&log.VersionHash,
			&log.AgentID,
			&log.Metadata,
			&log.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}
	return logs, rows.Err()
}

// DeleteOldLogs deletes audit logs older than the specified number of days
func (r *AuditRepository) DeleteOldLogs(days int) (int64, error) {
	query := `DELETE FROM audit_logs WHERE created_at < NOW() - INTERVAL '1 day' * $1`
	result, err := r.db.Exec(query, days)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old audit logs: %w", err)
	}
	deletedCount, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get deleted count: %w", err)
	}
	return deletedCount, nil
}

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

