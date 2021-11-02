package commits

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var (
	dir     = ""
	Command = &cobra.Command{
		Use:   "commits",
		Short: "Utility functions to deal with Git commits",
		Run:   system.Help,
	}
)

func init() {
	Command.PersistentFlags().StringVar(&dir, "git-dir", "./", "")
	Command.AddCommand(&cobra.Command{Use: "list", Run: list})
}
