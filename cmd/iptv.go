package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// iptvCmd represents the iptv command
var iptvCmd = &cobra.Command{
	Use:   "iptv",
	Short: "Select and filter your favored IPTV channels",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("iptv called")
	},
}

func init() {
	rootCmd.AddCommand(iptvCmd)
}
