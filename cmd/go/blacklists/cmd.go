package blacklists

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var (
	Command = &cobra.Command{
		Use: "blacklists",
		Run: system.Help,
	}
)

func init() {
	Command.AddCommand(&cobra.Command{Use: "analyze", Run: analyze})
}
