package nodes

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use:   "nodes",
	Short: "Utility functions to deal with Nodes",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "selectors", Run: selectors})
}
