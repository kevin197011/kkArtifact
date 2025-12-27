// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"net"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// getAgentIDFromRequest extracts agent identifier from the request
// Format: hostname-ip
func getAgentIDFromRequest(c *gin.Context) string {
	// Get client IP address
	clientIP := getClientIP(c)
	if clientIP == "" {
		clientIP = "unknown"
	}

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// Format: hostname-ip
	return hostname + "-" + clientIP
}

// getClientIP extracts the client IP address from the request
// It checks X-Forwarded-For, X-Real-IP headers, and RemoteAddr
func getClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header (for proxies/load balancers)
	xff := c.GetHeader("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can contain multiple IPs, get the first one
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			// Validate IP
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	// Check X-Real-IP header
	xri := c.GetHeader("X-Real-IP")
	if xri != "" {
		ip := strings.TrimSpace(xri)
		if net.ParseIP(ip) != nil {
			return ip
		}
	}

	// Fall back to RemoteAddr
	remoteAddr := c.Request.RemoteAddr
	if remoteAddr != "" {
		// RemoteAddr format: "ip:port"
		host, _, err := net.SplitHostPort(remoteAddr)
		if err == nil && host != "" {
			return host
		}
		// If SplitHostPort fails, try parsing as IP directly
		if net.ParseIP(remoteAddr) != nil {
			return remoteAddr
		}
	}

	return ""
}

