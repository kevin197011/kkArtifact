// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Client is the API client
type Client struct {
	serverURL  string
	token      string
	httpClient *http.Client
}

// New creates a new API client with optimized HTTP transport for high concurrency
func New(serverURL, token string) *Client {
	// Clean token: remove any existing "Bearer " prefix to avoid duplication
	cleanToken := strings.TrimSpace(token)
	if strings.HasPrefix(cleanToken, "Bearer ") {
		cleanToken = strings.TrimPrefix(cleanToken, "Bearer ")
		cleanToken = strings.TrimSpace(cleanToken)
	}

	// Configure HTTP transport with connection pooling for better performance
	transport := &http.Transport{
		MaxIdleConns:        500,              // Maximum idle connections across all hosts
		MaxIdleConnsPerHost: 200,              // Maximum idle connections per host (important for concurrent uploads)
		MaxConnsPerHost:     300,              // Maximum connections per host (prevents overwhelming server)
		IdleConnTimeout:     90 * time.Second, // How long idle connections are kept
		DisableCompression:  false,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 300 * time.Second, // For large file transfers
	}

	return &Client{
		serverURL: serverURL,
		token:     cleanToken, // Store cleaned token without "Bearer " prefix
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   600 * time.Second, // Total request timeout (10 minutes for large files)
		},
	}
}

// UploadInitRequest represents upload init request
type UploadInitRequest struct {
	Project   string `json:"project"`
	App       string `json:"app"`
	Version   string `json:"version"`
	FileCount int    `json:"file_count"`
}

// UploadInitResponse represents upload init response
type UploadInitResponse struct {
	UploadID string `json:"upload_id"`
}

// InitUpload initializes an upload session
func (c *Client) InitUpload(project, app, version string, fileCount int) (*UploadInitResponse, error) {
	req := UploadInitRequest{
		Project:   project,
		App:       app,
		Version:   version,
		FileCount: fileCount,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", c.serverURL+"/api/v1/upload/init", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	// Set Authorization header (token is already cleaned in New(), so we can safely add "Bearer ")
	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		// Output detailed request information for debugging
		fmt.Fprintf(os.Stderr, "\n=== Init Upload Request Details ===\n")
		fmt.Fprintf(os.Stderr, "URL: %s\n", c.serverURL+"/api/v1/upload/init")
		fmt.Fprintf(os.Stderr, "Method: %s\n", httpReq.Method)
		fmt.Fprintf(os.Stderr, "Request Headers:\n")
		for key, values := range httpReq.Header {
			for _, value := range values {
				// Print Authorization header in full for debugging (mask only the actual token value)
				if key == "Authorization" {
					// Print full header value to see if there's duplicate "Bearer"
					// Mask the actual token part (after "Bearer " or "Bearer Bearer ")
					prefix := ""
					tokenPart := value
					if strings.HasPrefix(value, "Bearer Bearer ") {
						prefix = "Bearer Bearer "
						tokenPart = value[len("Bearer Bearer "):]
					} else if strings.HasPrefix(value, "Bearer ") {
						prefix = "Bearer "
						tokenPart = value[len("Bearer "):]
					}
					if len(tokenPart) > 20 {
						fmt.Fprintf(os.Stderr, "  %s: %s%s...%s [FULL: %s]\n", key, prefix, tokenPart[:10], tokenPart[len(tokenPart)-10:], value)
					} else {
						fmt.Fprintf(os.Stderr, "  %s: %s%s [FULL: %s]\n", key, prefix, tokenPart, value)
					}
				} else {
					fmt.Fprintf(os.Stderr, "  %s: %s\n", key, value)
				}
			}
		}
		fmt.Fprintf(os.Stderr, "Response Status: %d %s\n", resp.StatusCode, resp.Status)
		fmt.Fprintf(os.Stderr, "Response Body: %s\n", string(body))
		fmt.Fprintf(os.Stderr, "====================================\n\n")

		return nil, fmt.Errorf("upload init failed with status %d: %s", resp.StatusCode, string(body))
	}

	var uploadResp UploadInitResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, err
	}

	return &uploadResp, nil
}

// UploadFile uploads a single file
func (c *Client) UploadFile(project, app, hash, filePath, localPath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add path field
	if err := writer.WriteField("path", filePath); err != nil {
		return err
	}

	// Add file
	part, err := writer.CreateFormFile("file", filepath.Base(localPath))
	if err != nil {
		return err
	}

	if _, err := io.Copy(part, file); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/v1/file/%s/%s/%s", c.serverURL, project, app, hash)
	httpReq, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	// Ensure token is always set before making request
	if c.token == "" {
		return fmt.Errorf("token is empty, cannot upload file")
	}
	// Set Authorization header (token is already cleaned in New(), so we can safely add "Bearer ")
	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		// Output detailed request information for debugging
		fmt.Fprintf(os.Stderr, "\n=== Upload Request Details ===\n")
		fmt.Fprintf(os.Stderr, "URL: %s\n", url)
		fmt.Fprintf(os.Stderr, "Method: %s\n", httpReq.Method)
		fmt.Fprintf(os.Stderr, "Request Headers:\n")
		for key, values := range httpReq.Header {
			for _, value := range values {
				// Print Authorization header in full for debugging (mask only the actual token value)
				if key == "Authorization" {
					// Print full header value to see if there's duplicate "Bearer"
					// Mask the actual token part (after "Bearer " or "Bearer Bearer ")
					prefix := ""
					tokenPart := value
					if strings.HasPrefix(value, "Bearer Bearer ") {
						prefix = "Bearer Bearer "
						tokenPart = value[len("Bearer Bearer "):]
					} else if strings.HasPrefix(value, "Bearer ") {
						prefix = "Bearer "
						tokenPart = value[len("Bearer "):]
					}
					if len(tokenPart) > 20 {
						fmt.Fprintf(os.Stderr, "  %s: %s%s...%s [FULL: %s]\n", key, prefix, tokenPart[:10], tokenPart[len(tokenPart)-10:], value)
					} else {
						fmt.Fprintf(os.Stderr, "  %s: %s%s [FULL: %s]\n", key, prefix, tokenPart, value)
					}
				} else {
					fmt.Fprintf(os.Stderr, "  %s: %s\n", key, value)
				}
			}
		}
		fmt.Fprintf(os.Stderr, "Response Status: %d %s\n", resp.StatusCode, resp.Status)
		fmt.Fprintf(os.Stderr, "Response Headers:\n")
		for key, values := range resp.Header {
			for _, value := range values {
				fmt.Fprintf(os.Stderr, "  %s: %s\n", key, value)
			}
		}
		fmt.Fprintf(os.Stderr, "Response Body: %s\n", string(body))
		fmt.Fprintf(os.Stderr, "==============================\n\n")

		if resp.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("upload failed with status %d (unauthorized): %s. Please check your token in the config file", resp.StatusCode, string(body))
		}
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// FinishUpload finishes the upload
func (c *Client) FinishUpload(req interface{}) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest("POST", c.serverURL+"/api/v1/upload/finish", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	// Set Authorization header (token is already cleaned in New(), so we can safely add "Bearer ")
	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		// Output detailed request information for debugging
		fmt.Fprintf(os.Stderr, "\n=== Finish Upload Request Details ===\n")
		fmt.Fprintf(os.Stderr, "URL: %s\n", c.serverURL+"/api/v1/upload/finish")
		fmt.Fprintf(os.Stderr, "Method: %s\n", httpReq.Method)
		fmt.Fprintf(os.Stderr, "Request Headers:\n")
		for key, values := range httpReq.Header {
			for _, value := range values {
				// Print Authorization header in full for debugging (mask only the actual token value)
				if key == "Authorization" {
					// Print full header value to see if there's duplicate "Bearer"
					// Mask the actual token part (after "Bearer " or "Bearer Bearer ")
					prefix := ""
					tokenPart := value
					if strings.HasPrefix(value, "Bearer Bearer ") {
						prefix = "Bearer Bearer "
						tokenPart = value[len("Bearer Bearer "):]
					} else if strings.HasPrefix(value, "Bearer ") {
						prefix = "Bearer "
						tokenPart = value[len("Bearer "):]
					}
					if len(tokenPart) > 20 {
						fmt.Fprintf(os.Stderr, "  %s: %s%s...%s [FULL: %s]\n", key, prefix, tokenPart[:10], tokenPart[len(tokenPart)-10:], value)
					} else {
						fmt.Fprintf(os.Stderr, "  %s: %s%s [FULL: %s]\n", key, prefix, tokenPart, value)
					}
				} else {
					fmt.Fprintf(os.Stderr, "  %s: %s\n", key, value)
				}
			}
		}
		fmt.Fprintf(os.Stderr, "Response Status: %d %s\n", resp.StatusCode, resp.Status)
		fmt.Fprintf(os.Stderr, "Response Body: %s\n", string(body))
		fmt.Fprintf(os.Stderr, "======================================\n\n")

		return fmt.Errorf("finish upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetManifest retrieves a manifest
func (c *Client) GetManifest(project, app, version string) (interface{}, error) {
	// Ensure token is set before making request
	if c.token == "" {
		return nil, fmt.Errorf("token is empty, cannot get manifest. Please check your config file (global: /etc/kkArtifact/config.yml or local: .kkartifact.yml)")
	}

	url := fmt.Sprintf("%s/api/v1/manifest/%s/%s/%s", c.serverURL, project, app, version)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read error response body for more details
		body, _ := io.ReadAll(resp.Body)
		errorMsg := fmt.Sprintf("get manifest failed with status %d", resp.StatusCode)
		if resp.StatusCode == http.StatusUnauthorized {
			errorMsg += " (unauthorized). Please check your token in the config file (global: /etc/kkArtifact/config.yml or local: .kkartifact.yml)"
		}
		if len(body) > 0 {
			errorMsg += fmt.Sprintf(": %s", string(body))
		}
		return nil, fmt.Errorf(errorMsg)
	}

	var manifest interface{}
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, err
	}

	return manifest, nil
}

// LatestVersionResponse represents the latest version response
type LatestVersionResponse struct {
	Project string `json:"project"`
	App     string `json:"app"`
	Version string `json:"version"`
}

// GetLatestVersion retrieves the latest published version for an app
func (c *Client) GetLatestVersion(project, app string) (*LatestVersionResponse, error) {
	// Ensure token is set before making request
	if c.token == "" {
		return nil, fmt.Errorf("token is empty, cannot get latest version. Please check your config file (global: /etc/kkArtifact/config.yml or local: .kkartifact.yml)")
	}

	url := fmt.Sprintf("%s/api/v1/projects/%s/apps/%s/latest", c.serverURL, project, app)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		errorMsg := fmt.Sprintf("get latest version failed with status %d", resp.StatusCode)
		if resp.StatusCode == http.StatusUnauthorized {
			errorMsg += " (unauthorized). Please check your token in the config file (global: /etc/kkArtifact/config.yml or local: .kkartifact.yml)"
		}
		if len(body) > 0 {
			errorMsg += fmt.Sprintf(": %s", string(body))
		}
		return nil, fmt.Errorf(errorMsg)
	}

	var latestResp LatestVersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&latestResp); err != nil {
		return nil, err
	}

	return &latestResp, nil
}

// AgentVersionInfo represents agent version information
type AgentVersionInfo struct {
	Version   string            `json:"version"`
	BuildTime string            `json:"build_time"`
	Binaries  []AgentBinaryInfo `json:"binaries"`
}

// AgentBinaryInfo represents a single agent binary
type AgentBinaryInfo struct {
	Platform string `json:"platform"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	URL      string `json:"url"`
}

// GetAgentVersionInfo retrieves agent version information from the server
func (c *Client) GetAgentVersionInfo() (*AgentVersionInfo, error) {
	url := fmt.Sprintf("%s/api/v1/downloads/agent/version", c.serverURL)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Agent version endpoint is public, no auth required
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get agent version info failed with status %d: %s", resp.StatusCode, string(body))
	}

	var versionInfo AgentVersionInfo
	if err := json.NewDecoder(resp.Body).Decode(&versionInfo); err != nil {
		return nil, err
	}

	return &versionInfo, nil
}

// DownloadAgentBinary downloads an agent binary file
func (c *Client) DownloadAgentBinary(filename string, destPath string) error {
	url := fmt.Sprintf("%s/api/v1/downloads/agent/%s", c.serverURL, filename)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Agent download endpoint is public, no auth required
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("download agent binary failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Create destination file
	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer file.Close()

	// Copy content
	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Make file executable on Unix systems
	if err := os.Chmod(destPath, 0755); err != nil {
		// Ignore chmod errors on Windows
	}

	return nil
}
