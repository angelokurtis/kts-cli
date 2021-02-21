package cmd

import (
	"github.com/spf13/cobra"
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Utility functions for tracking changes in any set of files using Git",
}

func init() {
	rootCmd.AddCommand(gitCmd)
}
