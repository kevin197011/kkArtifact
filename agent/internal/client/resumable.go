// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package client

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// CalculateFileHash calculates SHA256 hash of a file
func CalculateFileHash(filePath string) (string, int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	hash := sha256.New()
	size, err := io.Copy(hash, file)
	if err != nil {
		return "", 0, err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), size, nil
}

// CheckFileExistsAndMatches checks if a local file exists and matches the expected hash
// Returns (exists, matches, size, error)
func CheckFileExistsAndMatches(localPath, expectedHash string) (bool, bool, int64, error) {
	_, err := os.Stat(localPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, false, 0, nil
		}
		return false, false, 0, err
	}

	// File exists, check hash
	actualHash, size, err := CalculateFileHash(localPath)
	if err != nil {
		return true, false, size, err
	}

	matches := actualHash == expectedHash
	return true, matches, size, nil
}
