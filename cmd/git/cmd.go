package git

import (
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "git",
	Short: "git version-control utilities",
	Run:   common.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "sign-commits", Run: signCommits})
}
