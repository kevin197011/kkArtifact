// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// TokenRepository handles token database operations
type TokenRepository struct {
	db *DB
}

// NewTokenRepository creates a new token repository
func NewTokenRepository(db *DB) *TokenRepository {
	return &TokenRepository{db: db}
}

// Create creates a new token
func (r *TokenRepository) Create(tokenHash, name string, projectID, appID *int, permissions []string, expiresAt *time.Time) (*Token, error) {
	var token Token
	query := `INSERT INTO tokens (token_hash, name, project_id, app_id, permissions, expires_at)
	          VALUES ($1, $2, $3, $4, $5, $6)
	          RETURNING id, token_hash, name, project_id, app_id, permissions, expires_at, created_at`
	
	var nameNull sql.NullString
	if name != "" {
		nameNull = sql.NullString{String: name, Valid: true}
	}
	
	err := r.db.QueryRow(
		query,
		tokenHash,
		nameNull,
		toNullInt64(projectID),
		toNullInt64(appID),
		pq.Array(permissions),
		toNullTime(expiresAt),
	).Scan(
		&token.ID,
		&token.TokenHash,
		&token.Name,
		&token.ProjectID,
		&token.AppID,
		pq.Array(&token.Permissions),
		&token.ExpiresAt,
		&token.CreatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}
	return &token, nil
}

// FindByHash finds a token by its hash
func (r *TokenRepository) FindByHash(tokenHash string) (*Token, error) {
	var token Token
	query := `SELECT id, token_hash, name, project_id, app_id, permissions, expires_at, created_at
	          FROM tokens WHERE token_hash = $1`
	
	err := r.db.QueryRow(query, tokenHash).Scan(
		&token.ID,
		&token.TokenHash,
		&token.Name,
		&token.ProjectID,
		&token.AppID,
		pq.Array(&token.Permissions),
		&token.ExpiresAt,
		&token.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find token: %w", err)
	}
	return &token, nil
}

// List lists all tokens
func (r *TokenRepository) List() ([]*Token, error) {
	query := `SELECT id, token_hash, name, project_id, app_id, permissions, expires_at, created_at
	          FROM tokens ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tokens: %w", err)
	}
	defer rows.Close()

	var tokens []*Token
	for rows.Next() {
		var token Token
		if err := rows.Scan(
			&token.ID,
			&token.TokenHash,
			&token.Name,
			&token.ProjectID,
			&token.AppID,
			pq.Array(&token.Permissions),
			&token.ExpiresAt,
			&token.CreatedAt,
		); err != nil {
			return nil, err
		}
		tokens = append(tokens, &token)
	}
	return tokens, rows.Err()
}

// Revoke revokes a token by ID
func (r *TokenRepository) Revoke(id int) error {
	query := `DELETE FROM tokens WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// Helper functions
func toNullInt64(i *int) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(*i), Valid: true}
}

func toNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

