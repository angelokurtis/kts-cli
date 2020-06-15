package gcp

import (
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/angelokurtis/kts-cli/cmd/gcp/container_registry"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "gcp",
	Short: "Google Cloud Platform utilities",
	Run:   common.Help,
}

func init() {
	Command.AddCommand(container_registry.Command)
}
