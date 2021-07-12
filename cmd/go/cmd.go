package golang

import (
	"github.com/angelokurtis/kts-cli/cmd/go/packages"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "go",
	Short: "Go utilities",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(packages.Command)
}
