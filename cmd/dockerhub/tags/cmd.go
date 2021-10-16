package tags

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use: "tags",
	Run: system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "list", Run: list, Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires the repository as argument")
		}
		return nil
	}})
}
