// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package util

import (
	"os"
	"strings"
)

// IsDebugMode checks if debug mode is enabled
func IsDebugMode() bool {
	debug := strings.ToLower(os.Getenv("DEBUG"))
	return debug == "true" || debug == "1" || debug == "yes"
}

