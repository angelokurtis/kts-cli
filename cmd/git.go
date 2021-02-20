package cmd

import (
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/spf13/cobra"
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Utility functions for tracking changes in any set of files using Git",
	Run:   system.Help,
}

func init() {
	rootCmd.AddCommand(gitCmd)
}
