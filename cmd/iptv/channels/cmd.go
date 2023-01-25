package channels

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/m3u"
	"github.com/spf13/cobra"
	"path/filepath"
)

var (
	filename = ""
	Command  = &cobra.Command{
		Use: "channels",
		Run: system.Help,
	}
)

func init() {
	Command.PersistentFlags().StringVarP(&filename, "filename", "f", "", "that contains the configuration to apply")
	Command.AddCommand(&cobra.Command{Use: "list", Run: list})
	Command.AddCommand(&cobra.Command{Use: "edit", Run: edit})
}

func list(cmd *cobra.Command, args []string) {
	channels, err := m3u.ListChannels(filename)
	dieOnErr(err)

	for _, channel := range channels {
		fmt.Printf("%s [%s]\n", channel.Name, channel.Group())
	}
}

func edit(cmd *cobra.Command, args []string) {
	channels, err := m3u.ListChannels(filename)
	dieOnErr(err)

	channels, err = channels.SelectMany()
	dieOnErr(err)

	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]

	err = channels.Write(name + "[edited]" + ext)
	dieOnErr(err)
}

func dieOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
