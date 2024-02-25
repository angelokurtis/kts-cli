package services

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use:   "services",
	Short: "Utility function to use port forwarding to access applications in a cluster",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "forwarding", Run: forwarding})
}
