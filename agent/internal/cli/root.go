// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cli

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// Version is set at build time using -ldflags
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "kkartifact-agent",
	Short: "kkArtifact Agent - CLI tool for artifact push/pull operations",
	Long:  "kkArtifact Agent is a command-line tool for pushing and pulling artifacts from kkArtifact server",
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, show help
		cmd.Help()
	},
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add version flag
	rootCmd.Flags().BoolP("version", "v", false, "Show version information")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// Check if version flag is set
		if version, _ := cmd.Flags().GetBool("version"); version {
			showVersion()
			os.Exit(0)
		}
	}
	
	// Subcommands are registered in their respective files
}

// showVersion displays version information
func showVersion() {
	fmt.Printf("kkartifact-agent version %s\n", Version)
	fmt.Printf("Build time: %s\n", BuildTime)
	fmt.Printf("Git commit: %s\n", GitCommit)
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

