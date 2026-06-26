// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kk/kkartifact-server/internal/auth"
)

func TestCallerHasPermission_JWTAdmin(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("session_info", &auth.SessionInfo{IsAdmin: true})

	for _, perm := range []auth.Permission{auth.PermissionPush, auth.PermissionPull, auth.PermissionPromote} {
		if !callerHasPermission(c, perm) {
			t.Errorf("admin session should have %s permission", perm)
		}
	}
}

func TestCallerHasPermission_JWTNonAdmin(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("session_info", &auth.SessionInfo{IsAdmin: false})

	if !callerHasPermission(c, auth.PermissionPull) {
		t.Error("non-admin session should have pull permission")
	}
	if callerHasPermission(c, auth.PermissionPush) {
		t.Error("non-admin session should not have push permission")
	}
}

func TestCallerHasPermission_APIToken(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("token_info", &auth.TokenInfo{Permissions: []string{"push"}})

	if !callerHasPermission(c, auth.PermissionPush) {
		t.Error("token with push should have push permission")
	}
	if callerHasPermission(c, auth.PermissionPull) {
		t.Error("token with only push should not have pull permission")
	}
}
