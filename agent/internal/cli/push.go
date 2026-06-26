// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cli

import (
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:          "push [flags]",
	Short:        "Push artifacts to the server",
	Long:         "Push artifacts from a local directory. App and version can be inferred from --path when it ends with app/version (e.g. .../G02_agent_api/cd884eb...).",
	SilenceUsage: true,
	RunE:         runPush,
}

var (
	pushProject     string
	pushApp         string
	pushVersion     string
	pushPath        string
	pushConfig      string
	pushServerURL   string
	pushToken       string
	pushConcurrency int
	pushIgnore      []string
	pushPublish     bool
)

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.Flags().StringVar(&pushApp, "app", "", "App name (optional if inferrable from --path)")
	pushCmd.Flags().StringVar(&pushVersion, "version", "", "Version hash (optional if inferrable from --path)")
	pushCmd.Flags().StringVar(&pushPath, "path", ".", "Path to version directory (last two segments: app/version)")
	bindPushConfigFlags(pushCmd, &pushProject, &pushConfig, &pushServerURL, &pushToken, &pushConcurrency, &pushIgnore)
	pushCmd.Flags().BoolVar(&pushPublish, "publish", false, "Mark version as published after push (available via pull latest)")
}

func runPush(cmd *cobra.Command, args []string) error {
	params, err := resolvePushParams(
		pushProject, pushApp, pushVersion, pushPath,
		pushConfig, pushServerURL, pushToken,
		pushConcurrency, pushIgnore,
	)
	if err != nil {
		return err
	}
	params.publish = pushPublish

	return executePush(params)
}
