package terraform

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/cmd/terraform/commands"
	"github.com/angelokurtis/kts-cli/cmd/terraform/providers"
	"github.com/angelokurtis/kts-cli/cmd/terraform/resources"
	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use:   "terraform",
	Short: "Utilities for HashiCorp Terraform",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(providers.Command)
	Command.AddCommand(resources.Command)
	Command.AddCommand(commands.Command)
}
