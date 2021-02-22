package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// iptvGroupsCmd represents the iptv groups command
var iptvGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List IPTV groups",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("iptv called")
	},
}

func init() {
	iptvCmd.AddCommand(iptvGroupsCmd)
}
