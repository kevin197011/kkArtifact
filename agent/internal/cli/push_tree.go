// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cli

import (
	"fmt"

	"github.com/kk/kkartifact-agent/internal/client"
	"github.com/spf13/cobra"
)

var pushTreeCmd = &cobra.Command{
	Use:          "push-tree [root-path]",
	Short:        "Push all app/version directories under a root path",
	Long:         "Traverse root/{app}/{version}/ and push each version. Matches layout like /data/vcs/G02/tidb/G02_agent_api/<hash>/",
	SilenceUsage: true,
	RunE:         runPushTree,
}

var (
	pushTreeRoot         string
	pushTreeProject      string
	pushTreeConfig       string
	pushTreeServerURL    string
	pushTreeToken        string
	pushTreeConcurrency  int
	pushTreeIgnore       []string
	pushTreeDryRun       bool
	pushTreeSkipExisting bool
	pushTreePublish      bool
)

func init() {
	rootCmd.AddCommand(pushTreeCmd)

	pushTreeCmd.Flags().StringVar(&pushTreeRoot, "root", "", "Root directory containing app/version subdirectories")
	bindPushConfigFlags(pushTreeCmd, &pushTreeProject, &pushTreeConfig, &pushTreeServerURL, &pushTreeToken, &pushTreeConcurrency, &pushTreeIgnore)
	pushTreeCmd.Flags().BoolVar(&pushTreeDryRun, "dry-run", false, "Print planned pushes without uploading")
	pushTreeCmd.Flags().BoolVar(&pushTreeSkipExisting, "skip-existing", false, "Skip versions that already exist on the server")
	pushTreeCmd.Flags().BoolVar(&pushTreePublish, "publish", false, "Mark each version as published after push")
}

func runPushTree(cmd *cobra.Command, args []string) error {
	root := pushTreeRoot
	if len(args) > 0 {
		root = args[0]
	}
	if root == "" {
		return fmt.Errorf("root path is required (argument or --root)")
	}

	targets, err := listPushTreeTargets(root)
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		fmt.Println("No app/version directories found")
		return nil
	}

	baseCfg, project, err := loadPushConfig(
		pushTreeProject,
		pushTreeConfig, pushTreeServerURL, pushTreeToken,
		pushTreeConcurrency, pushTreeIgnore,
	)
	if err != nil {
		return err
	}

	apiClient, err := client.New(baseCfg.ServerURL, baseCfg.Token)
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	fmt.Printf("Found %d version(s) under %s (project: %s)\n", len(targets), root, project)

	var failed int
	for _, target := range targets {
		params := &pushParams{
			project:   project,
			app:       target.app,
			version:   target.version,
			path:      target.path,
			cfg:       baseCfg,
			publish:   pushTreePublish,
			apiClient: apiClient,
		}

		label := fmt.Sprintf("%s/%s", params.app, params.version)
		if pushTreeDryRun {
			fmt.Printf("[dry-run] would push %s/%s:%s from %s\n", params.project, params.app, params.version, params.path)
			continue
		}

		if pushTreeSkipExisting && versionExistsOnServer(apiClient, params.project, params.app, params.version) {
			fmt.Printf(">>> skip existing %s\n", label)
			continue
		}

		fmt.Printf(">>> push %s\n", label)
		if err := executePush(params); err != nil {
			fmt.Printf("ERROR: %s: %v\n", label, err)
			failed++
		}
		fmt.Println()
	}

	if failed > 0 {
		return fmt.Errorf("%d push(es) failed", failed)
	}
	return nil
}
