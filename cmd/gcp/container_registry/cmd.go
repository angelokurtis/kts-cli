package container_registry

import (
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/spf13/cobra"
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
