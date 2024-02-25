package iptv

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/cmd/iptv/channels"
	"github.com/angelokurtis/kts-cli/cmd/iptv/groups"
	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use:   "iptv",
	Short: "IPTV functions utilities",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(channels.Command)
	Command.AddCommand(groups.Command)
}
