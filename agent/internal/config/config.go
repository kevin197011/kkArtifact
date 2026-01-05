// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package config

import (
	"fmt"
	"os"
	"path/filepath"

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
	Concurrency    int      `yaml:"concurrency"` // Number of concurrent uploads/downloads (default: 50)
}

// GetGlobalConfigPath returns the path to the global configuration file
// System-wide global config: /etc/kkartifact/kkartifact.yml
func GetGlobalConfigPath() (string, error) {
	// System-wide global configuration
	systemConfigPath := "/etc/kkartifact/kkartifact.yml"
	return systemConfigPath, nil
}

// loadConfigFile loads a single configuration file
func loadConfigFile(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// mergeConfigs merges global config with local config
// Local config values override global config values
func mergeConfigs(global, local *Config) *Config {
	result := &Config{}

	// Start with global config
	if global != nil {
		result.ServerURL = global.ServerURL
		result.Token = global.Token
		result.Project = global.Project
		result.App = global.App
		result.Ignore = global.Ignore
		result.RetainVersions = global.RetainVersions
		result.Concurrency = global.Concurrency
	}

	// Override with local config (if present)
	if local != nil {
		if local.ServerURL != "" {
			result.ServerURL = local.ServerURL
		}
		if local.Token != "" {
			result.Token = local.Token
		}
		if local.Project != "" {
			result.Project = local.Project
		}
		if local.App != "" {
			result.App = local.App
		}
		// Ignore: if local config has ignore field (even if empty), use it
		// This allows clearing global ignore list by setting ignore: [] in local config
		if local.Ignore != nil {
			result.Ignore = local.Ignore
		}
		if local.RetainVersions != nil {
			result.RetainVersions = local.RetainVersions
		}
		if local.Concurrency > 0 {
			result.Concurrency = local.Concurrency
		}
	}

	return result
}

// Load loads configuration with priority: local config > global config
// If configPath is empty or ".kkartifact.yml", it will try to load from current directory
// Global config is loaded from /etc/kkartifact/kkartifact.yml
func Load(configPath string) (*Config, error) {
	var globalConfig *Config
	var localConfig *Config
	var err error

	// Try to load global config (ignore errors if it doesn't exist)
	globalConfigPath, err := GetGlobalConfigPath()
	if err == nil {
		if cfg, err := loadConfigFile(globalConfigPath); err == nil {
			globalConfig = cfg
		}
		// Ignore error if global config doesn't exist
	}

	// Load local config (required)
	// If configPath is empty or ".kkartifact.yml", use current directory
	if configPath == "" || configPath == ".kkartifact.yml" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}
		configPath = filepath.Join(wd, ".kkartifact.yml")
	}

	localConfig, err = loadConfigFile(configPath)
	if err != nil {
		// If local config doesn't exist and global config exists, use global config
		if globalConfig != nil {
			// Validate global config
			if globalConfig.ServerURL == "" {
				return nil, fmt.Errorf("server_url is required in global config")
			}
			if globalConfig.Token == "" {
				return nil, fmt.Errorf("token is required in global config")
			}
			// Set default concurrency if not specified
			if globalConfig.Concurrency <= 0 {
				globalConfig.Concurrency = 50
			}
			return globalConfig, nil
		}
		return nil, fmt.Errorf("failed to load config file %s: %w", configPath, err)
	}

	// Merge configs (local overrides global)
	mergedConfig := mergeConfigs(globalConfig, localConfig)

	// Validate required fields
	if mergedConfig.ServerURL == "" {
		return nil, fmt.Errorf("server_url is required")
	}
	if mergedConfig.Token == "" {
		return nil, fmt.Errorf("token is required")
	}

	// Set default concurrency if not specified or invalid
	if mergedConfig.Concurrency <= 0 {
		mergedConfig.Concurrency = 50 // Default to 50 concurrent operations
	}

	return mergedConfig, nil
}

