// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package scheduler

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/kk/kkartifact-server/internal/database"
	"github.com/kk/kkartifact-server/internal/util"
)

// AuditCleanupTask is a scheduled task for audit log cleanup
type AuditCleanupTask struct {
	db *database.DB
}

// NewAuditCleanupTask creates a new audit log cleanup task
func NewAuditCleanupTask(db *database.DB) *AuditCleanupTask {
	return &AuditCleanupTask{
		db: db,
	}
}

// Name returns the task name
func (t *AuditCleanupTask) Name() string {
	return "audit-log-cleanup"
}

// Run runs the audit log cleanup task
func (t *AuditCleanupTask) Run(ctx context.Context) error {
	// Get retention days from config
	configRepo := database.NewConfigRepository(t.db)
	retentionDaysStr, err := configRepo.Get("audit_log_retention_days")
	if err != nil {
		return fmt.Errorf("failed to get audit log retention days: %w", err)
	}

	retentionDays, err := strconv.Atoi(retentionDaysStr)
	if err != nil {
		return fmt.Errorf("invalid audit log retention days: %w", err)
	}

	// Delete old audit logs
	auditRepo := database.NewAuditRepository(t.db)
	deletedCount, err := auditRepo.DeleteOldLogs(retentionDays)
	if err != nil {
		return fmt.Errorf("failed to delete old audit logs: %w", err)
	}

	if deletedCount > 0 && util.IsDebugMode() {
		log.Printf("Deleted %d audit log entries older than %d days", deletedCount, retentionDays)
	}

	return nil
}

