package aws

import (
	"github.com/angelokurtis/kts-cli/cmd/aws/ecr"
	"github.com/angelokurtis/kts-cli/cmd/aws/profiles"
	"github.com/angelokurtis/kts-cli/cmd/aws/route53"
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
	Command.AddCommand(route53.Command)
	Command.AddCommand(profiles.Command)
}
