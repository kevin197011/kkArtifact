// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package auth

import (
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken()
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if len(token) == 0 {
		t.Error("Generated token is empty")
	}

	// Generate another token and verify they're different
	token2, err := GenerateToken()
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == token2 {
		t.Error("Generated tokens should be different")
	}
}

func TestHashToken(t *testing.T) {
	token := "test-token"
	hash, err := HashToken(token)
	if err != nil {
		t.Fatalf("Failed to hash token: %v", err)
	}

	if len(hash) == 0 {
		t.Error("Hashed token is empty")
	}

	// Hash should be different from original
	if hash == token {
		t.Error("Hash should be different from original token")
	}
}

func TestVerifyToken(t *testing.T) {
	token := "test-token"
	hash, err := HashToken(token)
	if err != nil {
		t.Fatalf("Failed to hash token: %v", err)
	}

	// Verify correct token
	if !VerifyToken(token, hash) {
		t.Error("Token verification failed for correct token")
	}

	// Verify incorrect token
	if VerifyToken("wrong-token", hash) {
		t.Error("Token verification should fail for incorrect token")
	}
}

func TestIsExpired(t *testing.T) {
	now := time.Now()

	// Non-expired token
	future := now.Add(1 * time.Hour)
	if IsExpired(&future) {
		t.Error("Future time should not be expired")
	}

	// Expired token
	past := now.Add(-1 * time.Hour)
	if !IsExpired(&past) {
		t.Error("Past time should be expired")
	}

	// Nil expiry (never expires)
	if IsExpired(nil) {
		t.Error("Nil expiry should not be expired")
	}
}

func TestHasPermission(t *testing.T) {
	tests := []struct {
		name        string
		permissions []string
		required    Permission
		want        bool
	}{
		{"has push permission", []string{"push", "pull"}, PermissionPush, true},
		{"missing push permission", []string{"pull"}, PermissionPush, false},
		{"admin has all permissions", []string{"admin"}, PermissionPush, true},
		{"admin has pull permission", []string{"admin"}, PermissionPull, true},
		{"no permissions", []string{}, PermissionPush, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasPermission(tt.permissions, tt.required)
			if got != tt.want {
				t.Errorf("HasPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}
