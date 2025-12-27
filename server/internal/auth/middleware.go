// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kk/kkartifact-server/internal/database"
	"github.com/lib/pq"
)

// TokenAuthenticator handles token authentication
type TokenAuthenticator struct {
	db *database.DB
}

// NewTokenAuthenticator creates a new token authenticator
func NewTokenAuthenticator(db *database.DB) *TokenAuthenticator {
	return &TokenAuthenticator{db: db}
}

// AuthenticateRequest authenticates a request and returns session info or token info
// Supports both JWT tokens (for Web UI sessions) and API tokens (for agents)
func (ta *TokenAuthenticator) AuthenticateRequest(c *gin.Context) (*SessionInfo, *TokenInfo, error) {
	// Get Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, nil, fmt.Errorf("missing authorization header")
	}

	// Check if it's a Bearer token
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, nil, fmt.Errorf("invalid authorization format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Try JWT token first (for Web UI sessions)
	jwtClaims, err := ValidateJWTToken(token)
	if err == nil && jwtClaims != nil {
		// Valid JWT token - return session info
		return &SessionInfo{
			UserID:   jwtClaims.UserID,
			Username: jwtClaims.Username,
			IsAdmin:  jwtClaims.IsAdmin,
		}, nil, nil
	}

	// If JWT validation fails, try API token (for agents)
	tokenInfo, err := ta.AuthenticateAPIToken(c.Request)
	if err == nil && tokenInfo != nil {
		// Valid API token - return token info
		return nil, tokenInfo, nil
	}

	// Both failed
	return nil, nil, fmt.Errorf("invalid token")
}

// AuthMiddleware returns a Gin middleware for authentication
// Supports both JWT tokens (Web UI) and API tokens (agents)
func (ta *TokenAuthenticator) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionInfo, tokenInfo, err := ta.AuthenticateRequest(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Store session info or token info in context
		if sessionInfo != nil {
			c.Set("session_info", sessionInfo)
		} else if tokenInfo != nil {
			c.Set("token_info", tokenInfo)
		}

		c.Next()
	}
}

// AuthenticateAPIToken authenticates an API token and returns token info
// This is for agent authentication (API tokens stored in database)
func (ta *TokenAuthenticator) AuthenticateAPIToken(req *http.Request) (*TokenInfo, error) {
	// Get Authorization header
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("missing authorization header")
	}

	// Check if it's a Bearer token
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, fmt.Errorf("invalid authorization format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Query database for all non-expired tokens and verify the token
	// Note: This is a simplified implementation. In production, consider:
	// 1. Adding a token lookup cache
	// 2. Using a token prefix index for faster lookup
	// 3. Limiting the number of tokens checked
	query := `SELECT id, project_id, app_id, permissions, expires_at, token_hash 
	          FROM tokens WHERE expires_at IS NULL OR expires_at > NOW()`
	rows, err := ta.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tokens: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var projectID, appID sql.NullInt64
		var permissions []string
		var expiresAt sql.NullTime
		var storedHash string

		if err := rows.Scan(&id, &projectID, &appID, pq.Array(&permissions), &expiresAt, &storedHash); err != nil {
			continue
		}

		// Verify token against stored hash
		if VerifyToken(token, storedHash) {
			var pID, aID *int
			if projectID.Valid {
				pid := int(projectID.Int64)
				pID = &pid
			}
			if appID.Valid {
				aid := int(appID.Int64)
				aID = &aid
			}

			return &TokenInfo{
				TokenID:     id,
				ProjectID:   pID,
				AppID:       aID,
				Permissions: permissions,
			}, nil
		}
	}

	return nil, fmt.Errorf("invalid token")
}
