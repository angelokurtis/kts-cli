package eks

import (
	"os"

	"github.com/gookit/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/aws"
)

var Command = &cobra.Command{
	Use: "eks",
	Run: system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "list", Run: list})
}

func list(cmd *cobra.Command, args []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Cluster", "Region"})
	table.SetBorder(false)

	for _, region := range aws.Regions {
		clusters, err := aws.ListEKSClusters(region)
		if err != nil {
			color.Yellow.Println("[WARN] " + err.Error())
		}

		for _, cluster := range clusters {
			table.Append([]string{cluster, region})
		}
	}

	table.Render()
}
