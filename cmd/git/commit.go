package git

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/git"
	"github.com/martinusso/inflect"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func commit(cmd *cobra.Command, args []string) {
	remoteBranch := "origin/HEAD"

	currentBranch, err := git.CurrentBranch()
	if err != nil {
		log.Fatal(err)
	}

	count, err := git.CountCommitsBetweenBranches(currentBranch, remoteBranch)
	if err != nil {
		log.Fatal(err)
	}

	message := fmt.Sprintf("Commit number %s", inflect.IntoWords(float64(count+1)))

	files, err := git.UncommittedFiles()
	if err != nil {
		log.Fatal(err)
	}

	current, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	paths := make([]string, 0, len(files))
	for _, file := range files {
		path, err := filepath.Rel(current, file)
		if err != nil {
			log.Fatal(err)
		}
		paths = append(paths, path)
	}

	paths, err = selectFiles(paths)
	if err != nil {
		log.Fatal(err)
	}

	if err = git.DoCommit(message, paths); err != nil {
		log.Fatal(err)
	}
}

func selectFiles(files []string) ([]string, error) {
	if files == nil || len(files) == 0 {
		return nil, errors.New("there's nothing available to commit")
	}

	var selects []string
	prompt := &survey.MultiSelect{
		Message: "Choose the files you want to commit:",
		Options: files,
	}

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(20), survey.WithKeepFilter(true))
	if err != nil {
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
