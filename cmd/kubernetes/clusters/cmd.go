package clusters

import (
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "clusters",
	Short: "Utility function to interact with all available Kubernetes clusters",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "list", Run: list})
	Command.AddCommand(&cobra.Command{Use: "config", Run: config})
}
