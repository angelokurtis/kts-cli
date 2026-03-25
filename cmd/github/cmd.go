package github

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/cmd/github/actions"
	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use:   "github",
	Short: "Interact with GitHub resources and automation",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(actions.Command)
}
