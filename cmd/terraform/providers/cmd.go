package providers

import (
	"fmt"
	"io/ioutil"
	log "log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/terraform"
)

var Command = &cobra.Command{
	Use:   "providers",
	Short: "Utility functions to deal with Providers in the Terraform Registry",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "import", Run: importCmd})
}

func importCmd(cmd *cobra.Command, args []string) {
	provider, err := terraform.SelectProvider()
	if err != nil {
		log.Error(err.Error())
		return
	}

	out, err := provider.Encode()
	if err != nil {
		log.Error(err.Error())
		return
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s.tf", provider.Name), out, os.ModePerm)
	if err != nil {
		log.Error(err.Error())
		return
	}
}
