package groups

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/m3u"
	"github.com/spf13/cobra"
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

	for _, group := range channels.Groups() {
		fmt.Println(group)
	}
}

func edit(cmd *cobra.Command, args []string) {

}

func dieOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
