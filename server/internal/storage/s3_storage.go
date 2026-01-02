// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3Storage implements Storage interface using S3-compatible storage
type S3Storage struct {
	client   *minio.Client
	bucket   string
	basePath string
}

// NewS3Storage creates a new S3 storage instance
func NewS3Storage(endpoint, accessKey, secretKey, bucket, basePath string, useSSL bool) (*S3Storage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %w", err)
	}

	return &S3Storage{
		client:   client,
		bucket:   bucket,
		basePath: basePath,
	}, nil
}

// buildPath constructs the full path including base path
func (s *S3Storage) buildPath(path string) string {
	cleanPath := strings.TrimPrefix(path, "/")
	if s.basePath != "" {
		return filepath.Join(s.basePath, cleanPath)
	}
	return cleanPath
}

// Put stores a file at the given path
func (s *S3Storage) Put(ctx context.Context, path string, reader io.Reader, size int64) error {
	objectPath := s.buildPath(path)
	_, err := s.client.PutObject(ctx, s.bucket, objectPath, reader, size, minio.PutObjectOptions{})
	return err
}

// Get retrieves a file from the given path
func (s *S3Storage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	objectPath := s.buildPath(path)
	obj, err := s.client.GetObject(ctx, s.bucket, objectPath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// Delete removes a file at the given path
func (s *S3Storage) Delete(ctx context.Context, path string) error {
	objectPath := s.buildPath(path)
	return s.client.RemoveObject(ctx, s.bucket, objectPath, minio.RemoveObjectOptions{})
}

// Exists checks if a path exists
func (s *S3Storage) Exists(ctx context.Context, path string) (bool, error) {
	objectPath := s.buildPath(path)
	_, err := s.client.StatObject(ctx, s.bucket, objectPath, minio.StatObjectOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// List lists files in a directory
func (s *S3Storage) List(ctx context.Context, prefix string) ([]string, error) {
	searchPrefix := s.buildPath(prefix)
	objectCh := s.client.ListObjects(context.Background(), s.bucket, minio.ListObjectsOptions{
		Prefix:    searchPrefix,
		Recursive: false,
	})

	var paths []string
	for obj := range objectCh {
		if obj.Err != nil {
			return nil, obj.Err
		}
		// Remove base path prefix from result
		relativePath := strings.TrimPrefix(obj.Key, s.basePath+"/")
		paths = append(paths, relativePath)
	}
	return paths, nil
}

// Stat returns metadata about a file
func (s *S3Storage) Stat(ctx context.Context, path string) (*FileInfo, error) {
	objectPath := s.buildPath(path)
	objInfo, err := s.client.StatObject(ctx, s.bucket, objectPath, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Path:    path,
		Size:    objInfo.Size,
		ModTime: objInfo.LastModified.Unix(),
		IsDir:   strings.HasSuffix(path, "/"),
	}, nil
}
