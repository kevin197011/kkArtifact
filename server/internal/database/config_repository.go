// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package database

import (
	"fmt"
)

// ConfigRepository handles configuration database operations
type ConfigRepository struct {
	db *DB
}

// NewConfigRepository creates a new config repository
func NewConfigRepository(db *DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

// Get gets a configuration value by key
func (r *ConfigRepository) Get(key string) (string, error) {
	var value string
	query := `SELECT value FROM config WHERE key = $1`
	err := r.db.QueryRow(query, key).Scan(&value)
	if err != nil {
		return "", fmt.Errorf("failed to get config: %w", err)
	}
	return value, nil
}

// Set sets a configuration value
func (r *ConfigRepository) Set(key, value string) error {
	query := `INSERT INTO config (key, value) VALUES ($1, $2)
	          ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = CURRENT_TIMESTAMP`
	_, err := r.db.Exec(query, key, value)
	return err
}

// GetAll gets all configuration entries
func (r *ConfigRepository) GetAll() (map[string]string, error) {
	query := `SELECT key, value FROM config`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all config: %w", err)
	}
	defer rows.Close()

	config := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		config[key] = value
	}
	return config, rows.Err()
}

