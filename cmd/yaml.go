package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// yamlCmd represents the yaml command
var yamlCmd = &cobra.Command{
	Use:   "yaml",
	Short: "Utility functions for dealing with YAML files",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("yaml called")
	},
}

func init() {
	rootCmd.AddCommand(yamlCmd)
}
