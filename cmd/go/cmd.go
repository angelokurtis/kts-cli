package golang

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/cmd/go/format"
	"github.com/angelokurtis/kts-cli/cmd/go/mod"
	"github.com/angelokurtis/kts-cli/cmd/go/packages"
	"github.com/angelokurtis/kts-cli/cmd/go/versions"
	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use:   "go",
	Short: "Go utilities",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(format.Command)
	Command.AddCommand(packages.NewCommand())
	Command.AddCommand(mod.Command)
	Command.AddCommand(versions.Command)
	Command.AddCommand(&cobra.Command{Use: "lint", Run: lint})
	// comment
	cleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "Remove comments from selected Go files",
		Run:   wrapWithErrorHandler(clean),
	}
	cleanCmd.Flags().BoolP("updated", "u", false, "Only select files with uncommitted Git changes")
	Command.AddCommand(cleanCmd)
}
