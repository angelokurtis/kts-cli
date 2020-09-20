package providers

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/terraformer"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "providers",
	Short: "Utility functions to deal with Providers in the Terraformer",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "list", Run: list})
}

func list(cmd *cobra.Command, args []string) {
	providers := terraformer.ListProviders()
	for _, provider := range providers {
		fmt.Println(provider)
	}
}
