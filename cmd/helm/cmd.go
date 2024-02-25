package helm

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/cmd/helm/revisions"
	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use: "helm",
	Run: system.Help,
}

func init() {
	Command.AddCommand(revisions.Command)
}
