// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package auth

// SessionInfo represents session information from JWT token
type SessionInfo struct {
	UserID   int
	Username string
	IsAdmin  bool
}

// TokenInfo represents API token information (for agent authentication)
type TokenInfo struct {
	TokenID     int
	ProjectID   *int
	AppID       *int
	Permissions []string
}

