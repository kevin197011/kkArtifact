// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

// init initializes JWT secret
func init() {
	// Try to get secret from environment variable
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Generate a random secret (will be different on each restart)
		// In production, should set JWT_SECRET environment variable
		secretBytes := make([]byte, 32)
		if _, err := rand.Read(secretBytes); err != nil {
			panic(fmt.Errorf("failed to generate JWT secret: %w", err))
		}
		secret = base64.URLEncoding.EncodeToString(secretBytes)
	}
	jwtSecret = []byte(secret)
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID   int      `json:"user_id"`
	Username string   `json:"username"`
	IsAdmin  bool     `json:"is_admin"`
	jwt.RegisteredClaims
}

// GenerateJWTToken generates a JWT token for a user
func GenerateJWTToken(userID int, username string, isAdmin bool) (string, error) {
	// Token expires in 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)
	
	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWTToken validates a JWT token and returns claims
func ValidateJWTToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

