package cmd

import (
	"github.com/angelokurtis/kts-cli/internal/system"

	"github.com/spf13/cobra"
)

// gitCloneCmd represents the git clone command
var gitCloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clones a repository into a newly created directory",
	Run:   system.Help,
}

func init() {
	gitCmd.AddCommand(gitCloneCmd)
}
