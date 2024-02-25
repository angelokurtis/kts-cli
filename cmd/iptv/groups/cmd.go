package groups

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/m3u"
)

var (
	filename = ""
	Command  = &cobra.Command{
		Use: "groups",
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

	for _, group := range channels.Groups().Items {
		fmt.Println(group)
	}
}

func edit(cmd *cobra.Command, args []string) {
	channels, err := m3u.ListChannels(filename)
	dieOnErr(err)

	groups, err := channels.Groups().SelectMany()
	dieOnErr(err)

	channels = channels.FilterByGroups(groups)

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
