package images

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/internal/system"
)

var tagged = false

var Command = &cobra.Command{
	Use: "images",
	Run: system.Help,
}

func init() {
	Command.PersistentFlags().BoolVarP(&tagged, "tagged", "t", false, "")
	Command.AddCommand(&cobra.Command{Use: "list", Run: list})
	Command.AddCommand(&cobra.Command{Use: "delete", Run: del})
}

func dieOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
