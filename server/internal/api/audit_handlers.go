// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kk/kkartifact-server/internal/database"
)

// AuditLogResponse represents an audit log entry in API response
type AuditLogResponse struct {
	ID          int     `json:"id"`
	Operation   string  `json:"operation"`
	ProjectID   *int    `json:"project_id,omitempty"`
	AppID       *int    `json:"app_id,omitempty"`
	VersionHash *string `json:"version_hash,omitempty"`
	AgentID     *string `json:"agent_id,omitempty"`
	Metadata    *string `json:"metadata,omitempty"`
	CreatedAt   string  `json:"created_at"` // RFC3339 format
}

// handleListAuditLogs lists audit logs
func (h *Handler) handleListAuditLogs(c *gin.Context) {
	projectID := getIntParam(c, "project_id", 0)
	appID := getIntParam(c, "app_id", 0)
	limit := getIntQuery(c, "limit", 50)
	offset := getIntQuery(c, "offset", 0)

	var projectIDPtr, appIDPtr *int
	if projectID > 0 {
		projectIDPtr = &projectID
	}
	if appID > 0 {
		appIDPtr = &appID
	}

	auditRepo := database.NewAuditRepository(h.db)
	logs, err := auditRepo.List(projectIDPtr, appIDPtr, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format with properly formatted dates
	responses := make([]AuditLogResponse, len(logs))
	for i, log := range logs {
		var projectID, appID *int
		if log.ProjectID.Valid {
			pid := int(log.ProjectID.Int64)
			projectID = &pid
		}
		if log.AppID.Valid {
			aid := int(log.AppID.Int64)
			appID = &aid
		}

		var versionHash, agentID, metadata *string
		if log.VersionHash.Valid {
			versionHash = &log.VersionHash.String
		}
		if log.AgentID.Valid {
			agentID = &log.AgentID.String
		}
		if log.Metadata.Valid {
			metadata = &log.Metadata.String
		}

		responses[i] = AuditLogResponse{
			ID:          log.ID,
			Operation:   log.Operation,
			ProjectID:   projectID,
			AppID:       appID,
			VersionHash: versionHash,
			AgentID:     agentID,
			Metadata:    metadata,
			CreatedAt:   log.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, responses)
}

// getIntParam is in helpers.go

