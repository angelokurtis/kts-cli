package profiles

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/aws"
)

// aws profiles list

var Command = &cobra.Command{
	Use:   "profiles",
	Short: "Utility functions to deal with Amazon profiles",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "list", Run: list})
}

func list(cmd *cobra.Command, args []string) {
	profiles, err := aws.ListProfiles()
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range profiles {
		fmt.Println(p)
	}
}
