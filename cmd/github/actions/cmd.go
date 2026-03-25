package actions

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use:   "actions",
	Short: "Manage and update GitHub Actions workflows",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "upgrade", Run: wrapWithErrorHandler(upgrade)})
}
