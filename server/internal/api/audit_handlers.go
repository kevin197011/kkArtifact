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
	ProjectName *string `json:"project_name,omitempty"`
	AppName     *string `json:"app_name,omitempty"`
	VersionHash *string `json:"version_hash,omitempty"`
	AgentID     *string `json:"agent_id,omitempty"`
	Metadata    *string `json:"metadata,omitempty"`
	CreatedAt   string  `json:"created_at"` // RFC3339 format
}

// AuditLogsListResponse represents the paginated audit logs API response
type AuditLogsListResponse struct {
	Data  []AuditLogResponse `json:"data"`
	Total int                `json:"total"`
}

// handleListAuditLogs godoc
// @Summary      List audit logs
// @Description  Get a list of audit logs with optional filtering by project and app. Returns paginated results with total count.
// @Tags         audit
// @Accept       json
// @Produce      json
// @Param        project_id  query     int  false  "Filter by project ID"
// @Param        app_id      query     int  false  "Filter by app ID"
// @Param        limit       query     int  false  "Limit number of results (default: 50)"
// @Param        offset      query     int  false  "Offset for pagination (default: 0)"
// @Success      200         {object}  AuditLogsListResponse
// @Failure      401         {object}  ErrorResponse
// @Failure      500         {object}  ErrorResponse
// @Security     Bearer
// @Router       /audit-logs [get]
func (h *Handler) handleListAuditLogs(c *gin.Context) {
	projectID := getIntQuery(c, "project_id", 0)
	appID := getIntQuery(c, "app_id", 0)
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
	
	// Get total count with same filters
	total, err := auditRepo.Count(projectIDPtr, appIDPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get paginated logs
	logs, err := auditRepo.List(projectIDPtr, appIDPtr, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format with properly formatted dates
	responses := make([]AuditLogResponse, len(logs))
	for i, log := range logs {
		var projectID, appID *int
		var projectName, appName *string
		
		if log.ProjectID.Valid {
			pid := int(log.ProjectID.Int64)
			projectID = &pid
			// Get project name by querying database directly
			var name string
			query := `SELECT name FROM projects WHERE id = $1`
			if err := h.db.QueryRow(query, pid).Scan(&name); err == nil {
				projectName = &name
			}
		}
		
		if log.AppID.Valid {
			aid := int(log.AppID.Int64)
			appID = &aid
			// Get app name by querying database directly
			var name string
			query := `SELECT name FROM apps WHERE id = $1`
			if err := h.db.QueryRow(query, aid).Scan(&name); err == nil {
				appName = &name
			}
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
			ProjectName: projectName,
			AppName:     appName,
			VersionHash: versionHash,
			AgentID:     agentID,
			Metadata:    metadata,
			CreatedAt:   log.CreatedAt.Format(time.RFC3339),
		}
	}

	// Return paginated response with total count
	c.JSON(http.StatusOK, AuditLogsListResponse{
		Data:  responses,
		Total: total,
	})
}

// getIntParam is in helpers.go

