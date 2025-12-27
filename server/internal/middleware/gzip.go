// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package middleware

import (
	"compress/gzip"
	"strings"

	"github.com/gin-gonic/gin"
)

// Gzip returns a gzip compression middleware
func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip compression for some content types
		if !shouldCompress(c) {
			c.Next()
			return
		}

		// Check if client accepts gzip
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		// Create gzip writer
		gz := gzip.NewWriter(c.Writer)
		defer gz.Close()

		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")
		c.Writer = &gzipWriter{ResponseWriter: c.Writer, Writer: gz}
		c.Next()
	}
}

type gzipWriter struct {
	gin.ResponseWriter
	Writer *gzip.Writer
}

func (w *gzipWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

func shouldCompress(c *gin.Context) bool {
	contentType := c.Writer.Header().Get("Content-Type")
	
	// Don't compress already compressed content
	if strings.Contains(contentType, "gzip") {
		return false
	}
	
	// Compress JSON and text responses
	if strings.HasPrefix(contentType, "application/json") ||
		strings.HasPrefix(contentType, "text/") {
		return true
	}
	
	return false
}

