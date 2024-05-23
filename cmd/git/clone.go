package git

import (
	"context"
	slog "log/slog"

	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/git"
	"github.com/angelokurtis/kts-cli/pkg/app/idea"
)

func clone(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	repo := args[0]

	dir, err := git.NewLocalDir(repo)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create local directory", slog.String("repo", repo), tint.Err(err))
		return
	}

	if !dir.Exist() {
		slog.InfoContext(ctx, "Cloning repository", slog.String("repo", repo))

		err = git.Clone(repo)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to clone repository", slog.String("repo", repo), tint.Err(err))
			return
		}

		slog.InfoContext(ctx, "Repository cloned successfully", slog.String("repo", repo))
	} else {
		slog.InfoContext(ctx, "Local repository found", slog.String("path", dir.Path()))
	}

	if open {
		slog.InfoContext(ctx, "Opening repository in IntelliJ IDEA", slog.String("path", dir.Path()))

		err = idea.Open(dir.Path())
		if err != nil {
			slog.ErrorContext(ctx, "Failed to open repository in IntelliJ IDEA", slog.String("path", dir.Path()), tint.Err(err))
			return
		}

		slog.InfoContext(ctx, "Repository opened in IntelliJ IDEA", slog.String("path", dir.Path()))
	}
}
