package images

import (
	"context"
	"log/slog"
	"sort"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func del(cmd *cobra.Command, args []string) {
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

	summaries, err = summaries.Select()
	dieOnErr(err)

	for _, summary := range summaries.wrapped {
		_, err = apiClient.ImageRemove(ctx, summary.wrapped.ID, image.RemoveOptions{Force: true, PruneChildren: true})
		if err != nil {
			slog.WarnContext(ctx, "Failed to remove image", slog.String("image-id", summary.wrapped.ID), slog.String("error", err.Error()))
			continue
		}

		slog.InfoContext(ctx, "Removed image", slog.String("image-id", summary.wrapped.ID))
	}
}
