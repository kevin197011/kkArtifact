// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// TokenScope represents the scope of a token
type TokenScope string

const (
	ScopeGlobal  TokenScope = "global"
	ScopeProject TokenScope = "project"
	ScopeApp     TokenScope = "app"
)

// Permission represents a permission type
type Permission string

const (
	PermissionPush    Permission = "push"
	PermissionPull    Permission = "pull"
	PermissionPromote Permission = "promote"
	PermissionAdmin   Permission = "admin"
)

// GenerateToken generates a new random token
func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// HashToken hashes a token using bcrypt
func HashToken(token string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash token: %w", err)
	}
	return string(hash), nil
}

// VerifyToken verifies a token against a hash
func VerifyToken(token, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
	return err == nil
}

// IsExpired checks if a token is expired
func IsExpired(expiresAt *time.Time) bool {
	if expiresAt == nil {
		return false
	}
	return time.Now().After(*expiresAt)
}

// HasPermission checks if permissions include the required permission
func HasPermission(permissions []string, required Permission) bool {
	for _, perm := range permissions {
		if perm == string(required) || perm == string(PermissionAdmin) {
			return true
		}
	}
	return false
}

