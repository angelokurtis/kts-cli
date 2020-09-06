package git

import (
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "git",
	Short: "git version-control utilities",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "sign-commits", Run: signCommits})
}
