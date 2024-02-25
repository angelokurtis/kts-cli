package gcp

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/cmd/gcp/container_registry"
	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use:   "gcp",
	Short: "Google Cloud Platform utilities",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(container_registry.Command)
}
