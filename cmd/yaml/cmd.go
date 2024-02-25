package yaml

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use:   "yaml",
	Short: "YAML functions utilities",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "split", Run: split, Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a YAML file argument")
		}
		return nil
	}})
	Command.AddCommand(&cobra.Command{Use: "join", Run: join, Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a YAML files path as argument")
		}
		return nil
	}})
}
