// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kk/kkartifact-server/internal/auth"
	"github.com/kk/kkartifact-server/internal/database"
)

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
// Token is a JWT token for session management (not stored in database)
type LoginResponse struct {
	Token string `json:"token"` // JWT token for session
	Name  string `json:"name"`  // Username
}

// handleLogin godoc
// @Summary      User login
// @Description  Login with username and password, returns JWT token for session management
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Login credentials"
// @Success      200      {object}  LoginResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Router       /login [post]
func (h *Handler) handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by username
	userRepo := database.NewUserRepository(h.db)
	user, err := userRepo.GetByUsername(req.Username)
	if err != nil {
		// Don't reveal whether username exists for security
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Verify password
	if !database.VerifyPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate JWT token for session (not stored in database)
	// For now, all users have admin permissions (TODO: implement roles)
	jwtToken, err := auth.GenerateJWTToken(user.ID, user.Username, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate session token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: jwtToken,
		Name:  user.Username,
	})
}

