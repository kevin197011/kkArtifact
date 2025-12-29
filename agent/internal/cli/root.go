// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kkartifact-agent",
	Short: "kkArtifact Agent - CLI tool for artifact push/pull operations",
	Long:  "kkArtifact Agent is a command-line tool for pushing and pulling artifacts from kkArtifact server",
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Subcommands are registered in their respective files
}

