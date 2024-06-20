package git

import (
	"context"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/git"
)

func stage(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	// Fetch uncommitted files
	files, err := git.UncommittedFiles()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch uncommitted files", tint.Err(err))
		os.Exit(1)
	}

	// Select files from the uncommitted files list
	selectedFiles, err := files.SelectFiles()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to select files from uncommitted files list", tint.Err(err))
		os.Exit(1)
	}

	// Determine unselected files by finding the difference
	unselectedFiles, _ := lo.Difference(files, selectedFiles)

	// Stage selected but unstaged files
	if err = stageSelectedFiles(selectedFiles); err != nil {
		slog.ErrorContext(ctx, "Failed to stage selected files", tint.Err(err))
		os.Exit(1)
	}

	// Unstage the files that were not selected but are currently staged
	if err = unstageUnselectedFiles(unselectedFiles); err != nil {
		slog.ErrorContext(ctx, "Failed to unstage unselected files", tint.Err(err))
		os.Exit(1)
	}
}

func stageSelectedFiles(files git.Files) error {
	unstaged, err := files.UnStagedFiles().RelativePaths()
	if err != nil {
		return err
	}

	return git.Stage(unstaged)
}

func unstageUnselectedFiles(files git.Files) error {
	staged, err := files.StagedFiles().RelativePaths()
	if err != nil {
		return err
	}

	return git.Unstage(staged)
}
