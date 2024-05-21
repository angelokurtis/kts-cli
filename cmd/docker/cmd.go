package docker

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/cmd/docker/images"
	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use: "docker",
	Run: system.Help,
}

func init() {
	Command.AddCommand(images.Command)
}
