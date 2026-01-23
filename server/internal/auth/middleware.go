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
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kk/kkartifact-server/internal/database"
	"github.com/lib/pq"
)

// cachedTokenInfo stores cached token information with expiration
type cachedTokenInfo struct {
	TokenInfo  *TokenInfo
	ExpiresAt  time.Time
	LastAccess time.Time
}

// TokenAuthenticator handles token authentication
type TokenAuthenticator struct {
	db            *database.DB
	tokenCache    sync.Map // map[string]*cachedTokenInfo - token hash -> cached info
	cacheTTL      time.Duration
	tokensLoadMux sync.RWMutex
	lastTokensLoad time.Time
	tokensCache   []tokenCacheEntry // Cached list of all tokens to avoid DB queries
	tokensCacheTTL time.Duration
}

// tokenCacheEntry stores token hash and info for fast lookup
type tokenCacheEntry struct {
	ID          int
	TokenHash   string
	ProjectID   *int
	AppID       *int
	Permissions []string
	ExpiresAt   *time.Time
}

// NewTokenAuthenticator creates a new token authenticator
func NewTokenAuthenticator(db *database.DB) *TokenAuthenticator {
	return &TokenAuthenticator{
		db:             db,
		cacheTTL:       5 * time.Minute,      // Cache validated tokens for 5 minutes
		tokensCacheTTL: 1 * time.Minute,      // Reload tokens list from DB every 1 minute
	}
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

// loadTokensFromDB loads all non-expired tokens from database (with caching)
func (ta *TokenAuthenticator) loadTokensFromDB() ([]tokenCacheEntry, error) {
	ta.tokensLoadMux.RLock()
	// Check if cache is still valid (within TTL)
	if time.Since(ta.lastTokensLoad) < ta.tokensCacheTTL && ta.tokensCache != nil {
		tokens := ta.tokensCache
		ta.tokensLoadMux.RUnlock()
		return tokens, nil
	}
	ta.tokensLoadMux.RUnlock()

	// Need to reload from database
	ta.tokensLoadMux.Lock()
	defer ta.tokensLoadMux.Unlock()

	// Double-check after acquiring write lock (another goroutine might have already loaded)
	if time.Since(ta.lastTokensLoad) < ta.tokensCacheTTL && ta.tokensCache != nil {
		return ta.tokensCache, nil
	}

	query := `SELECT id, project_id, app_id, permissions, expires_at, token_hash 
	          FROM tokens WHERE expires_at IS NULL OR expires_at > NOW()`
	rows, err := ta.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tokens: %w", err)
	}
	defer rows.Close()

	var tokens []tokenCacheEntry
	for rows.Next() {
		var id int
		var projectID, appID sql.NullInt64
		var permissions []string
		var expiresAt sql.NullTime
		var storedHash string

		if err := rows.Scan(&id, &projectID, &appID, pq.Array(&permissions), &expiresAt, &storedHash); err != nil {
			continue
		}

		var pID, aID *int
		if projectID.Valid {
			pid := int(projectID.Int64)
			pID = &pid
		}
		if appID.Valid {
			aid := int(appID.Int64)
			aID = &aid
		}

		var expAt *time.Time
		if expiresAt.Valid {
			expAt = &expiresAt.Time
		}

		tokens = append(tokens, tokenCacheEntry{
			ID:          id,
			TokenHash:   storedHash,
			ProjectID:   pID,
			AppID:       aID,
			Permissions: permissions,
			ExpiresAt:   expAt,
		})
	}

	ta.tokensCache = tokens
	ta.lastTokensLoad = time.Now()
	return tokens, nil
}

// AuthenticateAPIToken authenticates an API token and returns token info
// This is for agent authentication (API tokens stored in database)
// Uses caching to avoid database queries on every request
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
	// Trim whitespace from token to handle edge cases
	token = strings.TrimSpace(token)

	// Check cache first (using token as key)
	if cached, ok := ta.tokenCache.Load(token); ok {
		cachedInfo := cached.(*cachedTokenInfo)
		if time.Now().Before(cachedInfo.ExpiresAt) {
			// Cache hit - update last access time
			cachedInfo.LastAccess = time.Now()
			return cachedInfo.TokenInfo, nil
		}
		// Cache expired, remove it
		ta.tokenCache.Delete(token)
	}

	// Cache miss - need to verify token
	// Load tokens from database (with caching to avoid frequent DB queries)
	tokens, err := ta.loadTokensFromDB()
	if err != nil {
		return nil, err
	}

	// Verify token against stored hashes
	for _, entry := range tokens {
		// Verify token against stored hash
		if VerifyToken(token, entry.TokenHash) {
			tokenInfo := &TokenInfo{
				TokenID:     entry.ID,
				ProjectID:   entry.ProjectID,
				AppID:       entry.AppID,
				Permissions: entry.Permissions,
			}

			// Cache the validated token
			cachedInfo := &cachedTokenInfo{
				TokenInfo:  tokenInfo,
				ExpiresAt:  time.Now().Add(ta.cacheTTL),
				LastAccess: time.Now(),
			}
			ta.tokenCache.Store(token, cachedInfo)

			return tokenInfo, nil
		}
	}

	// Log authentication failure (without exposing full token)
	maskedToken := maskTokenForLogging(token)
	return nil, fmt.Errorf("invalid token (masked: %s)", maskedToken)
}

// maskTokenForLogging returns a masked version of a token for logging
// Shows first 5 and last 5 characters: "YVza5...b_rc="
func maskTokenForLogging(token string) string {
	if len(token) <= 10 {
		return strings.Repeat("*", len(token))
	}
	return token[:5] + "..." + token[len(token)-5:]
}

// InvalidateTokenCache invalidates the token cache (call this when tokens are created/deleted)
func (ta *TokenAuthenticator) InvalidateTokenCache() {
	ta.tokensLoadMux.Lock()
	defer ta.tokensLoadMux.Unlock()
	ta.lastTokensLoad = time.Time{} // Reset to force reload
	ta.tokensCache = nil
}
