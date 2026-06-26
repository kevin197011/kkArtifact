// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"database/sql"
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
// @Description  Get a list of audit logs with optional filtering by project, app, and operation. Returns paginated results with total count.
// @Tags         audit
// @Accept       json
// @Produce      json
// @Param        project_id  query     int     false  "Filter by project ID"
// @Param        app_id      query     int     false  "Filter by app ID"
// @Param        operation   query     string  false  "Filter by operation (push, pull, publish, unpublish, token_create, token_delete, version_delete, etc.)"
// @Param        limit       query     int     false  "Limit number of results (default: 50, max: 500)"
// @Param        offset      query     int     false  "Offset for pagination (default: 0)"
// @Success      200         {object}  AuditLogsListResponse
// @Failure      401         {object}  ErrorResponse
// @Failure      500         {object}  ErrorResponse
// @Security     Bearer
// @Router       /audit-logs [get]
func (h *Handler) handleListAuditLogs(c *gin.Context) {
	projectID := getIntQuery(c, "project_id", 0)
	appID := getIntQuery(c, "app_id", 0)
	operation := c.Query("operation")
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

	total, err := auditRepo.Count(projectIDPtr, appIDPtr, operation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logs, err := auditRepo.List(projectIDPtr, appIDPtr, operation, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]AuditLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = auditLogToResponse(log)
	}

	c.JSON(http.StatusOK, AuditLogsListResponse{
		Data:  responses,
		Total: total,
	})
}

func auditLogToResponse(log *database.AuditLogWithNames) AuditLogResponse {
	return AuditLogResponse{
		ID:          log.ID,
		Operation:   log.Operation,
		ProjectID:   nullInt64Ptr(log.ProjectID),
		AppID:       nullInt64Ptr(log.AppID),
		ProjectName: nullStringPtr(log.ProjectName),
		AppName:     nullStringPtr(log.AppName),
		VersionHash: nullStringPtr(log.VersionHash),
		AgentID:     nullStringPtr(log.AgentID),
		Metadata:    nullStringPtr(log.Metadata),
		CreatedAt:   log.CreatedAt.Format(time.RFC3339),
	}
}

func nullInt64Ptr(v sql.NullInt64) *int {
	if !v.Valid {
		return nil
	}
	n := int(v.Int64)
	return &n
}

func nullStringPtr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}
