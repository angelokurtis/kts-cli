package files

import (
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use:     "files",
	Aliases: []string{"file"},
	Run:     system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "list", Aliases: []string{"ls"}, Run: list})
	Command.AddCommand(&cobra.Command{Use: "remove", Aliases: []string{"rm"}, Run: remove})
}

func check(err error) {
	if err != nil {
		slog.Error(err.Error())
		return
	}
}
