package iptv

import (
	"github.com/angelokurtis/kts-cli/cmd/iptv/channels"
	"github.com/angelokurtis/kts-cli/cmd/iptv/groups"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/spf13/cobra"
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
