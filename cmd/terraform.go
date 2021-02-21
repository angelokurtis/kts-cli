package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// terraformCmd represents the terraform command
var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Utility functions for dealing with HashiCorp Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("terraform called")
	},
}

func init() {
	rootCmd.AddCommand(terraformCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// terraformCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// terraformCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
