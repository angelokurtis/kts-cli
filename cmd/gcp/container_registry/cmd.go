package container_registry

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use:   "container-registry",
	Short: "Utility function of Docker images on Google Cloud Platform",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "list", Run: list})
	Command.AddCommand(&cobra.Command{Use: "untag", Run: untag})
	Command.AddCommand(&cobra.Command{Use: "clean", Run: clean})
}
