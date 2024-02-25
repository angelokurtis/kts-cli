package route53

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/aws"
)

// aws route53 records

var Command = &cobra.Command{
	Use:   "route53",
	Short: "Utility functions to deal with Amazon Route 53",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "records", Run: records})
}

func records(cmd *cobra.Command, args []string) {
	profiles, err := aws.SelectProfiles()
	if err != nil {
		log.Fatal(err)
	}

	items := make([]*aws.ResourceRecordSet, 0, 0)

	for _, profile := range profiles {
		records, err := aws.ListAllRecordsByProfile(profile)
		if err != nil {
			log.Fatal(err)
		}

		items = append(items, records.Items...)
	}

	records := &aws.ResourceRecordSets{Items: items}
	records.RenderTable()
}
