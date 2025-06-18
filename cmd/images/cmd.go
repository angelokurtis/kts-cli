package images

import (
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var rm bool

var Command = &cobra.Command{
	Use:     "images",
	Aliases: []string{"file"},
	Run:     system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "dedupe", Run: dedupe})
}

func check(err error) {
	if err != nil {
		slog.Error(err.Error())
		return
	}
}
