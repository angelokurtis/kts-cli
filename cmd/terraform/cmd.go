package terraform

import (
	"github.com/angelokurtis/kts-cli/cmd/terraform/providers"
	"github.com/angelokurtis/kts-cli/cmd/terraform/resources"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "terraform",
	Short: "Utilities for HashiCorp Terraform",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(providers.Command)
	Command.AddCommand(resources.Command)
}
