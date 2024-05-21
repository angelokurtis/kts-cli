package images

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func list(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	dieOnErr(err)
	defer apiClient.Close()

	rawSummaries, err := apiClient.ImageList(ctx, image.ListOptions{})
	dieOnErr(err)

	sort.Slice(rawSummaries, func(i, j int) bool {
		t1 := rawSummaries[i].Size
		t2 := rawSummaries[j].Size

		return t1 < t2
	})

	summaries := wrapImageSummaries(rawSummaries)
	if tagged {
		summaries = summaries.FilterTagged()
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetBorder(false)
	table.SetHeader([]string{"ID", "TAGS", "SIZE"})

	for _, summary := range summaries.wrapped {
		table.Append([]string{
			summary.wrapped.ID,
			strings.Join(summary.wrapped.RepoTags, "\n"),
			byteCount(summary.wrapped.Size),
		})
	}

	table.Render()
}

func byteCount(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
