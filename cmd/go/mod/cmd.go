package mod

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use: "mod",
	Run: system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "upgrade", Run: wrapWithErrorHandler(upgrade)})
}
