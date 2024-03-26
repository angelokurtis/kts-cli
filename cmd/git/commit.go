package git

import (
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/martinusso/inflect"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/git"
)

func commit(cmd *cobra.Command, args []string) {
	count, err := git.CountCommitsByAuthor()
	if err != nil {
		log.Fatal(err)
	}

	author, err := git.GetUser()
	if err != nil {
		authors := lo.Keys(count)

		author, err = selectAuthor(authors)
		if err != nil {
			log.Fatal(err)
		}
	}

	total := count[author]

	message := fmt.Sprintf("Commit number %s", inflect.IntoWords(float64(total+1)))

	files, err := git.UncommittedFiles()
	if err != nil {
		log.Fatal(err)
	}

	selectedFiles, err := files.SelectFiles()
	if err != nil {
		log.Fatal(err)
	}

	paths, err := selectedFiles.RelativePaths()
	if err != nil {
		log.Fatal(err)
	}

	if err = git.DoCommit(message, paths); err != nil {
		log.Fatal(err)
	}
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
