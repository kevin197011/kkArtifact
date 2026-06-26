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

// List lists audit logs with optional filters (includes project/app names via JOIN)
func (r *AuditRepository) List(projectID, appID *int, operation string, limit, offset int) ([]*AuditLogWithNames, error) {
	query := `SELECT al.id, al.operation, al.project_id, al.app_id, al.version_hash, al.agent_id, al.metadata, al.created_at,
	                 p.name AS project_name, a.name AS app_name
	          FROM audit_logs al
	          LEFT JOIN projects p ON al.project_id = p.id
	          LEFT JOIN apps a ON al.app_id = a.id
	          WHERE ($1::int IS NULL OR al.project_id = $1)
	          AND ($2::int IS NULL OR al.app_id = $2)
	          AND ($3::text IS NULL OR $3 = '' OR al.operation = $3)
	          ORDER BY al.created_at DESC
	          LIMIT $4 OFFSET $5`

	projectIDVal, appIDVal := auditFilterIDs(projectID, appID)

	rows, err := r.db.Query(query, projectIDVal, appIDVal, operation, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit logs: %w", err)
	}
	defer rows.Close()

	var logs []*AuditLogWithNames
	for rows.Next() {
		var log AuditLogWithNames
		if err := rows.Scan(
			&log.ID,
			&log.Operation,
			&log.ProjectID,
			&log.AppID,
			&log.VersionHash,
			&log.AgentID,
			&log.Metadata,
			&log.CreatedAt,
			&log.ProjectName,
			&log.AppName,
		); err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}
	return logs, rows.Err()
}

func auditFilterIDs(projectID, appID *int) (projectIDVal, appIDVal interface{}) {
	if projectID != nil {
		projectIDVal = *projectID
	}
	if appID != nil {
		appIDVal = *appID
	}
	return projectIDVal, appIDVal
}

// Count counts audit logs with optional filters
func (r *AuditRepository) Count(projectID, appID *int, operation string) (int, error) {
	query := `SELECT COUNT(*) 
	          FROM audit_logs 
	          WHERE ($1::int IS NULL OR project_id = $1)
	          AND ($2::int IS NULL OR app_id = $2)
	          AND ($3::text IS NULL OR $3 = '' OR operation = $3)`

	projectIDVal, appIDVal := auditFilterIDs(projectID, appID)

	var count int
	err := r.db.QueryRow(query, projectIDVal, appIDVal, operation).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count audit logs: %w", err)
	}
	return count, nil
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
