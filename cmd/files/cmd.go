package files

import (
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use: "files",
	Run: system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "list", Run: list})
}

func check(err error) {
	if err != nil {
		slog.Error(err.Error())
		return
	}
}
