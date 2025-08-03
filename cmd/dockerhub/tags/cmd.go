package tags

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var semver string

var Command = &cobra.Command{
	Use: "tags",
	Run: system.Help,
}

func init() {
	listCommand := &cobra.Command{Use: "list", Run: list, Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires the repository as argument")
		}
		return nil
	}}
	listCommand.PersistentFlags().StringVar(&semver, "semver", "", "Filter tags using a semantic versioning constraint (e.g. ^1.2.0, >=2.0.0 <3.0.0)")
	Command.AddCommand(listCommand)
}
