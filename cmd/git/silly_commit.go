package git

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/huh"
	changeCase "github.com/ku/go-change-case"
	"github.com/lmittmann/tint"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/ordinal"
	"github.com/angelokurtis/kts-cli/pkg/app/git"
)

func sillyCommit(cmd *cobra.Command, args []string) {
	ctx := context.Background() // Assuming a context is available

	files, err := git.ListStagedFiles()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list staged files", tint.Err(err))
		return
	}

	if len(files) == 0 {
		slog.WarnContext(ctx, "No staged files to commit")
		return
	}

	commits, err := git.CountCommitsByAuthor()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to count commits by author", tint.Err(err))
		return
	}

	author, err := git.GetUser()
	if err != nil {
		slog.WarnContext(ctx, "Failed to get user, selecting author manually", tint.Err(err))

		authors := lo.Keys(commits)

		author, err = selectAuthor(authors)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to select author", tint.Err(err))
			return
		}
	}

	total := commits[author]
	slog.DebugContext(ctx, "Commit count for author", slog.String("author", author), slog.Int64("total_commits", total))

	number, err := ordinal.FromNumber(total + 1)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to write commit number in ordinal word form", tint.Err(err))
		return
	}

	msg := fmt.Sprintf("%s commit", changeCase.UcFirst(number))

	slog.DebugContext(ctx, "Staged files detected", slog.Int("file_count", len(files)))

	var sb strings.Builder
	for _, file := range files {
		sb.WriteString("\t" + file + "\n")
	}

	fmt.Printf("◇  The following %d files will be committed:\n%s\n", len(files), sb.String())

	var confirm bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Affirmative("Yes").
				Negative("No").
				Title("Are you sure?").
				Value(&confirm),
		),
	)

	if err = form.WithTheme(huh.ThemeDracula()).Run(); err != nil {
		slog.ErrorContext(ctx, "unable to confirm the commit message", tint.Err(err))
		return
	}

	if !confirm {
		slog.DebugContext(ctx, "Commit message not confirmed by user")
		return
	}

	if err = git.DoCommitStagedFiles(msg); err != nil {
		slog.ErrorContext(ctx, "Failed to commit staged files", tint.Err(err))
		return
	}

	slog.DebugContext(ctx, "Commit successful", slog.String("message", msg))
}

func selectAuthor(authors []string) (string, error) {
	if authors == nil || len(authors) == 0 {
		return "", errors.New("there's no author available")
	}

	var selected string

	prompt := &survey.Select{
		Message: "Choose the commit author:",
		Options: authors,
	}

	err := survey.AskOne(prompt, &selected, survey.WithPageSize(20), survey.WithKeepFilter(true))
	if err != nil {
		return "", errors.WithStack(err)
	}

	return selected, nil
}

func selectFiles(files []string) ([]string, error) {
	if files == nil || len(files) == 0 {
		return nil, errors.New("there's nothing available to commit")
	}

	defaults, err := git.ListStagedFiles()
	if err != nil {
		return nil, err
	}

	var selects []string

	prompt := &survey.MultiSelect{
		Message: "Choose the files you want to commit:",
		Options: files,
		Default: defaults,
	}

	if err = survey.AskOne(prompt, &selects, survey.WithPageSize(20), survey.WithKeepFilter(true)); err != nil {
		return nil, errors.WithStack(err)
	}

	return selects, nil
}

func selectBranch(branches []string) ([]string, error) {
	if branches == nil || len(branches) == 0 {
		return nil, errors.New("there isn't a remote branch")
	}

	var selects []string

	prompt := &survey.MultiSelect{
		Message: "Choose the remote branch:",
		Options: branches,
	}

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(20), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return selects, nil
}
