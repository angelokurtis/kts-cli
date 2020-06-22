package resources

import (
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "resources",
	Short: "Utility function to deal with Kubernetes API resources available on the server",
	Run:   common.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "export", Run: export})
}
