package services

import (
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "services",
	Short: "Utility function to use port forwarding to access applications in a cluster",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "forwarding", Run: forwarding})
}
