package cmd

import (
	"fmt"
	"github.com/andanhm/go-prettytime"
	"github.com/angelokurtis/kts-cli/pkg/app/git"
	"github.com/olekukonko/tablewriter"
	"os"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

var (
	path = "./"
	// gitTagsCmd represents the git tags list command
	gitTagsCmd = &cobra.Command{
		Use:   "tags",
		Short: "List the existing tags in Git project",
		RunE: func(_ *cobra.Command, _ []string) error {
			tags, err := git.ListTags(path)
			if err != nil {
				return err
			}

			sort.Slice(tags, func(i, j int) bool {
				return tags[i].Time.After(*tags[j].Time)
			})

			brazil, err := time.LoadLocation("America/Sao_Paulo")
			if err != nil {
				return err
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"TAG", "COMMIT", "CREATED"})
			table.SetRowLine(false)
			table.SetBorder(false)
			table.SetColWidth(50)
			for _, tag := range tags {
				t := fmt.Sprintf("%s (%s)", tag.Time.In(brazil).Format("02/01/2006 15:04"), prettytime.Format(*tag.Time))
				table.Append([]string{tag.Name, tag.CommitID, t})
			}
			table.Render()
			return nil
		},
	}
)

func init() {
	gitCmd.AddCommand(gitTagsCmd)
}
