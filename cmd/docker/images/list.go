package images

import (
	"context"
	"fmt"
	log "log/slog"
	"os"
	"sort"
	"strings"
	"time"

	prettytime "github.com/andanhm/go-prettytime"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var brazil *time.Location

func init() {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Error(err.Error())
		return
	}

	brazil = loc
}

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
	table.SetHeader([]string{"ID", "TAGS", "SIZE", "CREATED"})

	for _, summary := range summaries.wrapped {
		createdTime := time.Unix(summary.wrapped.Created, 0).In(brazil)
		formattedTime := fmt.Sprintf("%s (%s)", createdTime.Format("02/01/2006 15:04"), prettytime.Format(createdTime))

		table.Append([]string{
			summary.wrapped.ID,
			strings.Join(summary.wrapped.RepoTags, "\n"),
			byteCount(summary.wrapped.Size),
			formattedTime,
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
