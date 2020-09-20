package resources

import (
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/terraformer"
	"github.com/spf13/cobra"
)

var (
	provider = ""
	Command  = &cobra.Command{
		Use:   "resources",
		Short: "Utility functions to deal with Resources in the Terraformer",
		Run:   system.Help,
	}
)

func init() {
	Command.PersistentFlags().StringVarP(&provider, "provider", "p", "", "")
	Command.AddCommand(&cobra.Command{Use: "import", Run: importCmd})
}

func importCmd(cmd *cobra.Command, args []string) {
	if provider == "" {
		providers := terraformer.ListProviders()
		p, err := providers.SelectProvider()
		if err != nil {
			log.Fatal(err)
		}
		provider = p
	}
	resources, err := terraformer.ListResources(provider)
	if err != nil {
		log.Fatal(err)
	}
	resources, err = resources.SelectMany()
	if err != nil {
		log.Fatal(err)
	}
	err = resources.Import()
	if err != nil {
		log.Fatal(err)
	}
}
