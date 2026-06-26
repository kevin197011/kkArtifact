// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/kk/kkartifact-agent/internal/client"
	"github.com/kk/kkartifact-agent/internal/manifest"
)

// executePush uploads artifacts for a single project/app/version
func executePush(params *pushParams) error {
	startTime := time.Now()

	fmt.Printf("Generating manifest for %s/%s:%s from %s\n", params.project, params.app, params.version, params.path)

	m, err := manifest.Generate(params.project, params.app, params.version, params.path, params.cfg.Ignore)
	if err != nil {
		return fmt.Errorf("failed to generate manifest: %w", err)
	}

	fmt.Printf("Found %d files\n", len(m.Files))

	apiClient := params.apiClient
	if apiClient == nil {
		var err error
		apiClient, err = client.New(params.cfg.ServerURL, params.cfg.Token)
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
	}

	filesToUpload := m.Files
	if remote, err := apiClient.GetManifest(params.project, params.app, params.version); err == nil {
		var fullyUnchanged bool
		filesToUpload, fullyUnchanged = filterFilesToUpload(m, remote)
		if fullyUnchanged {
			fmt.Printf("Version %s/%s:%s already up to date, skipping upload\n", params.project, params.app, params.version)
			return maybePublishVersion(apiClient, params)
		}
		skipped := len(m.Files) - len(filesToUpload)
		if skipped > 0 {
			fmt.Printf("Skipping %d unchanged file(s)\n", skipped)
		}
	}

	if len(filesToUpload) == 0 {
		fmt.Println("No files to upload")
		return maybePublishVersion(apiClient, params)
	}

	fmt.Println("Initializing upload...")
	uploadResp, err := apiClient.InitUpload(params.project, params.app, params.version, len(filesToUpload))
	if err != nil {
		return fmt.Errorf("failed to initialize upload: %w", err)
	}
	fmt.Printf("Upload ID: %s\n", uploadResp.UploadID)

	fmt.Printf("Uploading %d files with concurrency: %d\n", len(filesToUpload), params.cfg.Concurrency)

	progressBar := NewProgressBar(len(filesToUpload))

	type uploadTask struct {
		file      manifest.ManifestFile
		localPath string
	}

	tasks := make(chan uploadTask, len(filesToUpload))
	errors := make(chan error, len(filesToUpload))

	for _, file := range filesToUpload {
		tasks <- uploadTask{
			file:      file,
			localPath: filepath.Join(params.path, file.Path),
		}
	}
	close(tasks)

	var wg sync.WaitGroup
	for i := 0; i < params.cfg.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				if err := apiClient.UploadFile(params.project, params.app, params.version, task.file.Path, task.localPath); err != nil {
					errors <- fmt.Errorf("failed to upload file %s: %w", task.file.Path, err)
					return
				}
				progressBar.Update(1)
			}
		}()
	}

	wg.Wait()
	progressBar.Finish()
	close(errors)

	for err := range errors {
		if err != nil {
			return err
		}
	}

	fmt.Println("Finalizing upload...")
	finishReq := map[string]interface{}{
		"project":  params.project,
		"app":      params.app,
		"version":  params.version,
		"manifest": m,
	}

	if err := apiClient.FinishUpload(finishReq); err != nil {
		return fmt.Errorf("failed to finish upload: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("Successfully pushed %s/%s:%s\n", params.project, params.app, params.version)
	fmt.Printf("Total time: %v\n", duration.Round(time.Second))
	return maybePublishVersion(apiClient, params)
}

// maybePublishVersion publishes the version when --publish is set
func maybePublishVersion(apiClient *client.Client, params *pushParams) error {
	if !params.publish {
		return nil
	}
	fmt.Printf("Publishing %s/%s:%s...\n", params.project, params.app, params.version)
	if err := apiClient.Publish(params.project, params.app, params.version); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}
	fmt.Println("Version published")
	return nil
}

// listPushTreeTargets lists app/version directories under root for push-tree
func listPushTreeTargets(root string) ([]pushTreeTarget, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(absRoot)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("root path is not a directory: %s", absRoot)
	}

	var targets []pushTreeTarget
	appEntries, err := os.ReadDir(absRoot)
	if err != nil {
		return nil, err
	}

	for _, appEntry := range appEntries {
		if !appEntry.IsDir() || strings.HasPrefix(appEntry.Name(), ".") {
			continue
		}
		appName := appEntry.Name()
		appPath := filepath.Join(absRoot, appName)

		versionEntries, err := os.ReadDir(appPath)
		if err != nil {
			return nil, err
		}

		for _, versionEntry := range versionEntries {
			if !versionEntry.IsDir() || strings.HasPrefix(versionEntry.Name(), ".") {
				continue
			}
			targets = append(targets, pushTreeTarget{
				app:     appName,
				version: versionEntry.Name(),
				path:    filepath.Join(appPath, versionEntry.Name()),
			})
		}
	}

	return targets, nil
}
