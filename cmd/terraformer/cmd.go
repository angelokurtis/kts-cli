package terraformer

import (
	"github.com/angelokurtis/kts-cli/cmd/terraformer/providers"
	"github.com/angelokurtis/kts-cli/cmd/terraformer/resources"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "terraformer",
	Short: "Utilities to generate terraform files from existing infrastructure",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(providers.Command)
	Command.AddCommand(resources.Command)
}
