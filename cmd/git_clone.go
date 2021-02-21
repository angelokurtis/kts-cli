package cmd

import (
	"github.com/angelokurtis/kts-cli/internal/system"

	"github.com/spf13/cobra"
)

// gitCloneCmd represents the git clone command
var gitCloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "A brief description of your command",
	Run:   system.Help,
}

func init() {
	gitCmd.AddCommand(gitCloneCmd)
}
