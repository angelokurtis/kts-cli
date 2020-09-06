package gcp

import (
	"github.com/angelokurtis/kts-cli/cmd/gcp/container_registry"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "gcp",
	Short: "Google Cloud Platform utilities",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(container_registry.Command)
}
