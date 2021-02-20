package merge_requests

import (
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/spf13/cobra"
)

var (
	sortUpdated = false
	username    = ""
	Command     = &cobra.Command{
		Use:   "merge-requests",
		Short: "Utility function of projects on GitLab Merge Requests",
		Run:   system.Help,
	}
)

func init() {
	Command.PersistentFlags().StringVar(&username, "username", "", "")

	listCmd := &cobra.Command{Use: "list", Run: list}
	listCmd.PersistentFlags().BoolVar(&sortUpdated, "sort-updated", false, "")
	Command.AddCommand(listCmd)
	Command.AddCommand(&cobra.Command{Use: "assign", Run: assign})
}
