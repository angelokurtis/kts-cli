package aws

import (
	"github.com/angelokurtis/kts-cli/cmd/aws/ecr"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "aws",
	Short: "Utilities AWS environment",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(ecr.Command)
}
