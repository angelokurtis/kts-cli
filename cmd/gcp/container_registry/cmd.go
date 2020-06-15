package container_registry

import (
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "container-registry",
	Short: "Utility function of Docker images on Google Cloud Platform",
	Run: common.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "list", Run: list})
	Command.AddCommand(&cobra.Command{Use: "clean", Run: clean})
}
