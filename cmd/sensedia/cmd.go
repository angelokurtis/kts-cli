package sensedia

import (
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/angelokurtis/kts-cli/cmd/sensedia/servicemesh"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "sensedia",
	Short: "Sensedia utilities",
	Run:   common.Help,
}

func init() {
	Command.AddCommand(servicemesh.Command)
}
