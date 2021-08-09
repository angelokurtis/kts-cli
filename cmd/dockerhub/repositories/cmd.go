package repositories

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use: "repositories",
	Run: system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "list", Run: list})
}
