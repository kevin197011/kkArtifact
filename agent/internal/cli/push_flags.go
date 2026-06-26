// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cli

import "github.com/spf13/cobra"

// bindPushConfigFlags registers shared push/push-tree configuration flags.
func bindPushConfigFlags(
	cmd *cobra.Command,
	project, config, serverURL, token *string,
	concurrency *int,
	ignore *[]string,
) {
	cmd.Flags().StringVar(project, "project", "", "Project name (or set project in config)")
	cmd.Flags().StringVar(config, "config", ".kkartifact.yml", "Config file path")
	cmd.Flags().StringVar(serverURL, "server-url", "", "Server URL (overrides config file)")
	cmd.Flags().StringVar(token, "token", "", "Authentication token (overrides config file)")
	cmd.Flags().IntVar(concurrency, "concurrency", 0, "Number of concurrent uploads (overrides config file)")
	cmd.Flags().StringArrayVar(ignore, "ignore", []string{}, "Ignore patterns")
}
