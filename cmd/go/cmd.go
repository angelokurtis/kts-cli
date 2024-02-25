package golang

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/cmd/go/packages"
	"github.com/angelokurtis/kts-cli/cmd/go/versions"
	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use:   "go",
	Short: "Go utilities",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(packages.Command)
	Command.AddCommand(versions.Command)
	Command.AddCommand(&cobra.Command{Use: "lint", Run: lint})
}
