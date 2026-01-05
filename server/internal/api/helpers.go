// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Helper functions for API handlers

// getIntQuery gets an integer query parameter
func getIntQuery(c *gin.Context, key string, defaultValue int) int {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	if intValue < 1 {
		return defaultValue
	}
	if intValue > 500 {
		return 500
	}
	return intValue
}

// getIntParam gets an integer path parameter
func getIntParam(c *gin.Context, key string, defaultValue int) int {
	value := c.Param(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

