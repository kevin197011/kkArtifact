// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the agent configuration
type Config struct {
	ServerURL      string   `yaml:"server_url"`
	Token          string   `yaml:"token"`
	Project        string   `yaml:"project"`
	App            string   `yaml:"app"`
	Ignore         []string `yaml:"ignore,omitempty"`
	RetainVersions *int     `yaml:"retain_versions,omitempty"`
	Concurrency    int      `yaml:"concurrency"` // Number of concurrent uploads/downloads (default: 8)
}

// Load loads configuration from a YAML file
func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate required fields
	if config.ServerURL == "" {
		return nil, fmt.Errorf("server_url is required")
	}
	if config.Token == "" {
		return nil, fmt.Errorf("token is required")
	}

	// Set default concurrency if not specified or invalid
	if config.Concurrency <= 0 {
		config.Concurrency = 8 // Default to 8 concurrent operations
	}

	return &config, nil
}

