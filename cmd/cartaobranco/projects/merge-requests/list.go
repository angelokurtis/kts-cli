package merge_requests

import (
	"github.com/andanhm/go-prettytime"
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/gitlab"
	"github.com/enescakir/emoji"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"strconv"
	"strings"
)

// cartaobranco projects merge-requests list --username=tiago.angelo
func list(cmd *cobra.Command, args []string) {
	mrs, err := gitlab.SearchMergeRequestsByUser(username)
	if err != nil {
		log.Fatal(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"", "ID", "Description", "Assignees", "Votes", "Created", "Updated", "URL"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnSeparator("")
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetColWidth(50)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	if sortUpdated {
		sort.Slice(mrs, func(i, j int) bool {
			return mrs[i].UpdatedAt.Before(*mrs[j].UpdatedAt)
		})
	}
	for _, mr := range mrs {
		assignees := make([]string, 0, len(mr.Assignees))
		for _, assignee := range mr.Assignees {
			assignees = append(assignees, assignee.Username)
		}
		status := emoji.GreenCircle.String()

		if mr.MergeStatus != "can_be_merged" {
			status = emoji.RedCircle.String()
		} else {
			p, err := gitlab.PipelineByMergeRequest(mr)
			if err != nil {
				log.Fatal(err)
			}
			if p != nil && p.Status != "success" {
				status = emoji.YellowCircle.String()
			}
		}
		votes := ""
		if mr.Upvotes > 0 {
			votes = strconv.Itoa(mr.Upvotes) + emoji.ThumbsUp.String()
		}
		if mr.Downvotes > 0 {
			if len(votes) > 0 {
				votes = votes + " | "
			}
			votes = votes + strconv.Itoa(mr.Downvotes) + emoji.ThumbsDown.String()
		}
		table.Append([]string{status, strconv.Itoa(mr.IID), mr.SourceBranch, strings.Join(assignees, ", "), votes, prettytime.Format(*mr.CreatedAt), prettytime.Format(*mr.UpdatedAt), mr.WebURL})
	}

	table.Render()
}
