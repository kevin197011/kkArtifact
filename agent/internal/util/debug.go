// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package util

import (
	"os"
	"strings"
)

// IsDebugMode returns true when agent debug output is explicitly enabled.
//
// Supported values:
// - KKARTIFACT_AGENT_DEBUG=true|1|yes
// - DEBUG=true|1|yes (fallback)
func IsDebugMode() bool {
	if v := strings.ToLower(strings.TrimSpace(os.Getenv("KKARTIFACT_AGENT_DEBUG"))); v != "" {
		return v == "true" || v == "1" || v == "yes"
	}
	v := strings.ToLower(strings.TrimSpace(os.Getenv("DEBUG")))
	return v == "true" || v == "1" || v == "yes"
}

