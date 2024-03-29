package ingresses

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

// kube ingresses list
var Command = &cobra.Command{
	Use:   "ingresses",
	Short: "Utility functions to manages external access to the services in a cluster",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "list", Run: list})
	Command.AddCommand(&cobra.Command{Use: "hosts-mapping", Run: hostsMapping})
}
