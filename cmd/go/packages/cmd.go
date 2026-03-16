package packages

import (
	"github.com/spf13/cobra"
)

var (
	internal bool
	tests    bool
	tags     []string
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "packages [directory]",
		Short: "List Go package imports",
		Run:   packages,
	}
	cmd.Flags().BoolVarP(&internal, "internal", "i", false, "Show only internal imports")
	cmd.Flags().BoolVarP(&tests, "tests", "t", false, "Include test packages")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Build tags to pass to go list")

	return cmd
}
