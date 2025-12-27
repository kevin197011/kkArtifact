// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package client

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadFile downloads a file from the server with resume support
// If expectedHash is provided and local file matches, skip download
// If local file exists but hash doesn't match, resume from current position
func (c *Client) DownloadFile(project, app, version, filePath, localPath, expectedHash string, expectedSize int64) error {
	// Check if file already exists and matches hash
	if expectedHash != "" {
		exists, matches, size, err := CheckFileExistsAndMatches(localPath, expectedHash)
		if err != nil {
			return fmt.Errorf("failed to check local file: %w", err)
		}
		if exists && matches {
			// File exists and hash matches, skip download
			return nil
		}
		if exists && size > 0 && size < expectedSize {
			// File exists but incomplete, resume download
			return c.resumeDownload(project, app, version, filePath, localPath, size, expectedSize)
		}
		if exists {
			// File exists but hash doesn't match, remove and re-download
			if err := os.Remove(localPath); err != nil {
				return fmt.Errorf("failed to remove corrupted file: %w", err)
			}
		}
	}

	// Full download
	return c.fullDownload(project, app, version, filePath, localPath)
}

// fullDownload performs a full file download
func (c *Client) fullDownload(project, app, version, filePath, localPath string) error {
	url := fmt.Sprintf("%s/api/v1/file/%s/%s/%s?path=%s", c.serverURL, project, app, version, filePath)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Set Authorization header (token is already cleaned in New(), so we can safely add "Bearer ")
	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Create directory if needed
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return err
	}

	// Create or truncate file
	file, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy content
	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}

	return nil
}

// resumeDownload resumes a download from a specific byte position using HTTP Range request
func (c *Client) resumeDownload(project, app, version, filePath, localPath string, startByte, expectedSize int64) error {
	url := fmt.Sprintf("%s/api/v1/file/%s/%s/%s?path=%s", c.serverURL, project, app, version, filePath)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Add Range header for partial content
	// Format: bytes=start-end (end is inclusive, -1 means to end of file)
	if expectedSize > 0 && startByte < expectedSize {
		rangeHeader := fmt.Sprintf("bytes=%d-%d", startByte, expectedSize-1)
		httpReq.Header.Set("Range", rangeHeader)
	} else {
		rangeHeader := fmt.Sprintf("bytes=%d-", startByte)
		httpReq.Header.Set("Range", rangeHeader)
	}
	// Set Authorization header (token is already cleaned in New(), so we can safely add "Bearer ")
	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if server supports Range requests
	if resp.StatusCode == http.StatusPartialContent {
		// Open file in append mode
		file, err := os.OpenFile(localPath, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		// Copy remaining content
		if _, err := io.Copy(file, resp.Body); err != nil {
			return err
		}
		return nil
	} else if resp.StatusCode == http.StatusOK {
		// Server doesn't support Range requests, fallback to full download
		return c.fullDownload(project, app, version, filePath, localPath)
	} else {
		return fmt.Errorf("resume download failed with status %d", resp.StatusCode)
	}
}

// CheckFileExists checks if a file exists on the server
func (c *Client) CheckFileExists(project, app, version, filePath string) (bool, error) {
	url := fmt.Sprintf("%s/api/v1/file/%s/%s/%s?path=%s", c.serverURL, project, app, version, filePath)

	httpReq, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false, err
	}

	// Set Authorization header (token is already cleaned in New(), so we can safely add "Bearer ")
	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}
