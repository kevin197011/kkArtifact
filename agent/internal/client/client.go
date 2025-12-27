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
	"net/http"
	"os"
	"path/filepath"
)

// Client is the API client
type Client struct {
	serverURL string
	token     string
	httpClient *http.Client
}

// New creates a new API client
func New(serverURL, token string) *Client {
	return &Client{
		serverURL:  serverURL,
		token:      token,
		httpClient: &http.Client{},
	}
}

// UploadInitRequest represents upload init request
type UploadInitRequest struct {
	Project   string   `json:"project"`
	App       string   `json:"app"`
	Version   string   `json:"version"`
	FileCount int      `json:"file_count"`
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
	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upload init failed with status %d", resp.StatusCode)
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
	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
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
	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("finish upload failed with status %d", resp.StatusCode)
	}

	return nil
}

// GetManifest retrieves a manifest
func (c *Client) GetManifest(project, app, version string) (interface{}, error) {
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
		return nil, fmt.Errorf("get manifest failed with status %d", resp.StatusCode)
	}

	var manifest interface{}
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, err
	}

	return manifest, nil
}

