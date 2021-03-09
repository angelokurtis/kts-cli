package cmd

import (
	"github.com/spf13/cobra"
)

// dockerCmd represents the docker command
var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Utility functions for Docker containers",
}

func init() {
	rootCmd.AddCommand(dockerCmd)
}
