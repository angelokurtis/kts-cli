package packages

import (
	"github.com/spf13/cobra"
)

var internal bool

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "packages [directory]",
		Short: "List Go package imports",
		Run:   packages,
	}
	cmd.Flags().BoolVarP(&internal, "internal", "i", false, "Show only internal imports")

	return cmd
}
