package servicemesh

import (
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "service-mesh",
	Short: "Sensedia ServiceMesh utilities",
	Run:   common.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "login", Run: login})
}
