package git

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/cmd/git/commits"
	"github.com/angelokurtis/kts-cli/cmd/git/tags"
	"github.com/angelokurtis/kts-cli/internal/system"
)

var (
	open    = false
	Command = &cobra.Command{
		Use:   "git",
		Short: "git version-control utilities",
		Run:   system.Help,
	}
)

func init() {
	Command.AddCommand(tags.Command)
	Command.AddCommand(&cobra.Command{Use: "sign-commits", Run: signCommits})

	cloneCommand := &cobra.Command{Use: "clone", Run: clone}
	cloneCommand.PersistentFlags().BoolVar(&open, "open", false, "")
	Command.AddCommand(cloneCommand)
	Command.AddCommand(commits.Command)

	Command.AddCommand(&cobra.Command{Use: "commit", Run: commit})
	Command.AddCommand(&cobra.Command{Use: "stage", Run: stage})
}
