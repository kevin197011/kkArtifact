// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kk/kkartifact-server/internal/auth"
	"github.com/kk/kkartifact-server/internal/database"
)

// CreateTokenRequest represents a token creation request
type CreateTokenRequest struct {
	Name        string   `json:"name"`
	ProjectID   *int     `json:"project_id,omitempty"`
	AppID       *int     `json:"app_id,omitempty"`
	Permissions []string `json:"permissions"`
	ExpiresAt   *string  `json:"expires_at,omitempty"` // ISO 8601 format
}

// TokenResponse represents a token response
type TokenResponse struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Token       string   `json:"token,omitempty"` // Only returned on creation
	ProjectID   *int     `json:"project_id,omitempty"`
	AppID       *int     `json:"app_id,omitempty"`
	Permissions []string `json:"permissions"`
	ExpiresAt   *string  `json:"expires_at,omitempty"`
	CreatedAt   string   `json:"created_at"`
}

// handleCreateToken creates a new token
func (h *Handler) handleCreateToken(c *gin.Context) {
	var req CreateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate token
	token, err := auth.GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Hash token
	tokenHash, err := auth.HashToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash token"})
		return
	}

	// Parse expires_at if provided
	var expiresAt *time.Time
	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		parsed, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expires_at format"})
			return
		}
		expiresAt = &parsed
	}

	// Default permissions if not provided
	permissions := req.Permissions
	if len(permissions) == 0 {
		permissions = []string{"pull", "push", "promote"}
	}

	// Create token in database
	tokenRepo := database.NewTokenRepository(h.db)
	createdToken, err := tokenRepo.Create(
		tokenHash,
		req.Name,
		req.ProjectID,
		req.AppID,
		permissions,
		expiresAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Format response
	var expiresAtStr *string
	if createdToken.ExpiresAt.Valid {
		formatted := createdToken.ExpiresAt.Time.Format(time.RFC3339)
		expiresAtStr = &formatted
	}

	var name string
	if createdToken.Name.Valid {
		name = createdToken.Name.String
	}

	var projectID, appID *int
	if createdToken.ProjectID.Valid {
		pid := int(createdToken.ProjectID.Int64)
		projectID = &pid
	}
	if createdToken.AppID.Valid {
		aid := int(createdToken.AppID.Int64)
		appID = &aid
	}

	response := TokenResponse{
		ID:          createdToken.ID,
		Name:        name,
		Token:       token, // Return the plain token only once
		ProjectID:   projectID,
		AppID:       appID,
		Permissions: createdToken.Permissions,
		ExpiresAt:   expiresAtStr,
		CreatedAt:   createdToken.CreatedAt.Format(time.RFC3339),
	}

	// Record audit log for token creation
	auditRepo := database.NewAuditRepository(h.db)
	agentID := getAgentIDFromRequest(c)
	metadata := map[string]interface{}{
		"token_name": name,
		"permissions": permissions,
	}
	if projectID != nil {
		metadata["project_id"] = *projectID
	}
	if appID != nil {
		metadata["app_id"] = *appID
	}
	if expiresAtStr != nil {
		metadata["expires_at"] = *expiresAtStr
	}
	_ = auditRepo.Create("token_create", projectID, appID, "", agentID, metadata)

	c.JSON(http.StatusOK, response)
}

// handleListTokens lists all tokens
func (h *Handler) handleListTokens(c *gin.Context) {
	tokenRepo := database.NewTokenRepository(h.db)
	tokens, err := tokenRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format (without token values)
	responses := make([]TokenResponse, len(tokens))
	for i, token := range tokens {
		var expiresAtStr *string
		if token.ExpiresAt.Valid {
			formatted := token.ExpiresAt.Time.Format(time.RFC3339)
			expiresAtStr = &formatted
		}

		var name string
		if token.Name.Valid {
			name = token.Name.String
		}

		var projectID, appID *int
		if token.ProjectID.Valid {
			pid := int(token.ProjectID.Int64)
			projectID = &pid
		}
		if token.AppID.Valid {
			aid := int(token.AppID.Int64)
			appID = &aid
		}

		responses[i] = TokenResponse{
			ID:          token.ID,
			Name:        name,
			Token:       "", // Empty - token value is never returned in list (security)
			ProjectID:   projectID,
			AppID:       appID,
			Permissions: token.Permissions,
			ExpiresAt:   expiresAtStr,
			CreatedAt:   token.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, responses)
}

// handleDeleteToken deletes a token
func (h *Handler) handleDeleteToken(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token id"})
		return
	}

	tokenRepo := database.NewTokenRepository(h.db)
	
	// Get token info before deleting for audit log
	// List all tokens and find the one to delete
	tokens, err := tokenRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	var tokenToDelete *database.Token
	for _, token := range tokens {
		if token.ID == id {
			tokenToDelete = token
			break
		}
	}
	
	if tokenToDelete == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
		return
	}

	if err := tokenRepo.Revoke(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Record audit log for token deletion
	auditRepo := database.NewAuditRepository(h.db)
	agentID := getAgentIDFromRequest(c)
	var projectID, appID *int
	if tokenToDelete.ProjectID.Valid {
		pid := int(tokenToDelete.ProjectID.Int64)
		projectID = &pid
	}
	if tokenToDelete.AppID.Valid {
		aid := int(tokenToDelete.AppID.Int64)
		appID = &aid
	}
	metadata := map[string]interface{}{
		"token_id": id,
	}
	if tokenToDelete.Name.Valid {
		metadata["token_name"] = tokenToDelete.Name.String
	}
	_ = auditRepo.Create("token_delete", projectID, appID, "", agentID, metadata)

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
