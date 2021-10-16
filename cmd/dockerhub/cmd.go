package dockerhub

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/cmd/dockerhub/repositories"
	"github.com/angelokurtis/kts-cli/cmd/dockerhub/tags"
	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use:   "dockerhub",
	Short: "Docker Hub utilities",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(repositories.Command)
	Command.AddCommand(tags.Command)
}
