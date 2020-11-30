package tags

import (
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/spf13/cobra"
)

var (
	dir     = ""
	Command = &cobra.Command{
		Use:   "tags",
		Short: "Utility functions to deal with Git Tagging",
		Run:   system.Help,
	}
)

func init() {
	Command.PersistentFlags().StringVar(&dir, "git-dir", "./", "")
	Command.AddCommand(&cobra.Command{Use: "list", Run: list})
}
