// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error response
// @Description Error response structure
type ErrorResponse struct {
	Error   string `json:"error" example:"bad_request"`
	Message string `json:"message,omitempty" example:"Invalid request parameters"`
	Code    string `json:"code,omitempty" example:"INVALID_INPUT"`
}

// handleError handles API errors consistently
func handleError(c *gin.Context, statusCode int, err error, message ...string) {
	response := ErrorResponse{
		Error: err.Error(),
	}

	if len(message) > 0 {
		response.Message = message[0]
	}

	c.JSON(statusCode, response)
}

// handleNotFound handles 404 errors
func handleNotFound(c *gin.Context, resource string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Error:   "not_found",
		Message: resource + " not found",
	})
}

// handleBadRequest handles 400 errors
func handleBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error:   "bad_request",
		Message: message,
	})
}

// handleInternalError handles 500 errors
func handleInternalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error:   "internal_error",
		Message: "An internal error occurred",
	})
}

