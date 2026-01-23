// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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
// Unix/Linux/macOS: Tries /etc/kkArtifact/config.yml first (with capital A), then falls back to /etc/kkartifact/kkartifact.yml
// Windows: Uses C:\ProgramData\kkArtifact\config.yml
func GetGlobalConfigPath() (string, error) {
	if runtime.GOOS == "windows" {
		// Windows: Use ProgramData directory
		programData := os.Getenv("ProgramData")
		if programData == "" {
			programData = "C:\\ProgramData"
		}
		configPath := filepath.Join(programData, "kkArtifact", "config.yml")
		return configPath, nil
	}
	
	// Unix-like systems: Try primary path first: /etc/kkArtifact/config.yml (with capital A)
	primaryPath := "/etc/kkArtifact/config.yml"
	if _, err := os.Stat(primaryPath); err == nil {
		return primaryPath, nil
	}
	
	// Fallback to legacy path: /etc/kkartifact/kkartifact.yml (lowercase, backward compatibility)
	fallbackPath := "/etc/kkartifact/kkartifact.yml"
	return fallbackPath, nil
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

	// Trim whitespace from token and server_url to handle YAML formatting issues
	// YAML parser should handle inline comments correctly, but we'll trim whitespace
	config.Token = strings.TrimSpace(config.Token)
	config.ServerURL = strings.TrimSpace(config.ServerURL)
	
	// Additional safety: Remove any trailing comment-like text that might have been parsed incorrectly
	// This handles edge cases where YAML parser might include comment text in the value
	// Only remove if it looks like a comment pattern (starts with # after whitespace)
	if idx := strings.Index(config.Token, " #"); idx >= 0 {
		// Check if the part after # looks like a comment (common comment phrases)
		afterHash := strings.TrimSpace(config.Token[idx+2:])
		if strings.HasPrefix(afterHash, "Uncomment") || 
		   strings.HasPrefix(afterHash, "set your token") ||
		   strings.HasPrefix(afterHash, "YOUR_TOKEN") {
			// This looks like a comment, remove it
			config.Token = strings.TrimSpace(config.Token[:idx])
		}
	}

	return &config, nil
}

// mergeIgnorePatterns merges ignore patterns from multiple sources
// Priority: global → local → command-line
// Removes duplicates, keeping the last occurrence (command-line patterns appear last)
func mergeIgnorePatterns(global, local, commandLine []string) []string {
	// Combine all patterns in priority order
	allPatterns := make([]string, 0)
	
	// Add global patterns first
	if global != nil {
		allPatterns = append(allPatterns, global...)
	}
	
	// Add local patterns
	if local != nil {
		allPatterns = append(allPatterns, local...)
	}
	
	// Add command-line patterns last
	if commandLine != nil {
		allPatterns = append(allPatterns, commandLine...)
	}
	
	// Remove duplicates, keeping the last occurrence
	seen := make(map[string]bool)
	result := make([]string, 0)
	
	// Iterate in reverse to keep the last occurrence of each pattern
	for i := len(allPatterns) - 1; i >= 0; i-- {
		pattern := allPatterns[i]
		if !seen[pattern] {
			seen[pattern] = true
			result = append([]string{pattern}, result...)
		}
	}
	
	return result
}

// mergeConfigs merges global config with local config
// Local config values override global config values
// This is the legacy function for backward compatibility
func mergeConfigs(global, local *Config) *Config {
	return mergeConfigsWithOverrides(global, local, nil)
}

// Overrides represents command-line overrides for configuration
type Overrides struct {
	ServerURL   string
	Token       string
	Project     string
	App         string
	Ignore      []string
	Concurrency int // 0 means not set
}

// mergeConfigsWithOverrides merges global config, local config, and command-line overrides
// Priority: global config → local config → command-line overrides
func mergeConfigsWithOverrides(global, local *Config, overrides *Overrides) *Config {
	result := &Config{}

	// Start with global config
	if global != nil {
		result.ServerURL = strings.TrimSpace(global.ServerURL)
		result.Token = strings.TrimSpace(global.Token)
		result.Project = global.Project
		result.App = global.App
		result.Ignore = global.Ignore
		result.RetainVersions = global.RetainVersions
		result.Concurrency = global.Concurrency
	}

	// Override with local config (if present)
	// Only override if local config has non-empty values (empty strings don't override)
	if local != nil {
		if local.ServerURL != "" {
			result.ServerURL = strings.TrimSpace(local.ServerURL)
		}
		if local.Token != "" {
			// Only override if local token is non-empty (preserve global token if local is empty)
			result.Token = strings.TrimSpace(local.Token)
		}
		if local.Project != "" {
			result.Project = local.Project
		}
		if local.App != "" {
			result.App = local.App
		}
		if local.RetainVersions != nil {
			result.RetainVersions = local.RetainVersions
		}
		if local.Concurrency > 0 {
			result.Concurrency = local.Concurrency
		}
		// Note: ignore patterns are merged separately below
	}

	// Apply command-line overrides (highest priority)
	if overrides != nil {
		if overrides.ServerURL != "" {
			result.ServerURL = strings.TrimSpace(overrides.ServerURL)
		}
		if overrides.Token != "" {
			result.Token = strings.TrimSpace(overrides.Token)
		}
		if overrides.Project != "" {
			result.Project = overrides.Project
		}
		if overrides.App != "" {
			result.App = overrides.App
		}
		if overrides.Concurrency > 0 {
			result.Concurrency = overrides.Concurrency
		}
	}

	// Merge ignore patterns: global → local → command-line
	var globalIgnore, localIgnore, cmdIgnore []string
	if global != nil {
		globalIgnore = global.Ignore
	}
	if local != nil && local.Ignore != nil {
		localIgnore = local.Ignore
	}
	if overrides != nil && overrides.Ignore != nil {
		cmdIgnore = overrides.Ignore
	}
	
	// If command-line has ignore patterns, merge all sources
	if cmdIgnore != nil {
		result.Ignore = mergeIgnorePatterns(globalIgnore, localIgnore, cmdIgnore)
	} else if localIgnore != nil {
		// If local config has ignore (even if empty), use it (allows clearing global ignore)
		result.Ignore = mergeIgnorePatterns(globalIgnore, localIgnore, nil)
	} else {
		// Use global ignore only
		result.Ignore = globalIgnore
	}

	return result
}

// Load loads configuration with priority: global config → local config → command-line overrides
// If configPath is empty or ".kkartifact.yml", it will try to load from current directory
// Global config is loaded from /etc/kkArtifact/config.yml (or /etc/kkartifact/kkartifact.yml as fallback)
// If overrides is nil, Load behaves the same as before (backward compatible)
func Load(configPath string, overrides *Overrides) (*Config, error) {
	var globalConfig *Config
	var localConfig *Config
	var err error

	// Try to load global config (ignore errors if it doesn't exist)
	globalConfigPath, err := GetGlobalConfigPath()
	if err == nil {
		if runtime.GOOS == "windows" {
			// Windows: Use the path directly
			if cfg, err := loadConfigFile(globalConfigPath); err == nil {
				globalConfig = cfg
			}
			// Ignore error if global config doesn't exist
		} else {
			// Unix-like: Try primary path first
			primaryPath := "/etc/kkArtifact/config.yml"
			if cfg, err := loadConfigFile(primaryPath); err == nil {
				globalConfig = cfg
			} else {
				// Fallback to legacy path
				if cfg, err := loadConfigFile(globalConfigPath); err == nil {
					globalConfig = cfg
				}
				// Ignore error if global config doesn't exist
			}
		}
	}

	// Load local config (required if no global config or overrides)
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
			// Apply overrides if provided
			mergedConfig := mergeConfigsWithOverrides(globalConfig, nil, overrides)
			
			// Validate merged config
			if mergedConfig.ServerURL == "" {
				return nil, fmt.Errorf("server_url is required in global config")
			}
			if mergedConfig.Token == "" {
				return nil, fmt.Errorf("token is required in global config")
			}
			// Set default concurrency if not specified
			if mergedConfig.Concurrency <= 0 {
				mergedConfig.Concurrency = 50
			}
			return mergedConfig, nil
		}
		// If overrides provide required fields, we can proceed without config files
		if overrides != nil && overrides.ServerURL != "" && overrides.Token != "" {
			mergedConfig := mergeConfigsWithOverrides(nil, nil, overrides)
			if mergedConfig.Concurrency <= 0 {
				mergedConfig.Concurrency = 50
			}
			return mergedConfig, nil
		}
		return nil, fmt.Errorf("failed to load config file %s: %w", configPath, err)
	}

	// Merge configs: global → local → command-line overrides
	mergedConfig := mergeConfigsWithOverrides(globalConfig, localConfig, overrides)

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

