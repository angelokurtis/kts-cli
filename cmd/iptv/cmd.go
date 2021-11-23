package iptv

import (
	"os"

	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/m3u"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "iptv",
	Short: "IPTV functions utilities",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "channels", Run: channels})
}

// iptv channels
func channels(cmd *cobra.Command, args []string) {
	dirname, err := os.UserHomeDir()
	dieOnErr(err)

	channels, err := m3u.ListChannels(dirname + "/Downloads/tv_channels.m3u")
	dieOnErr(err)

	channels, err = channels.SelectMany()
	dieOnErr(err)

	err = channels.Write(dirname + "/Downloads/globo.m3u")
	dieOnErr(err)
}

func dieOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
