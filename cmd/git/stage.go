package git

import (
	log "log/slog"

	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/git"
)

func stage(cmd *cobra.Command, args []string) {
	// Fetch uncommitted files
	files, err := git.UncommittedFiles()
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Select files from the uncommitted files list
	selectedFiles, err := files.SelectFiles()
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Determine unselected files by finding the difference
	unselectedFiles, _ := lo.Difference(files, selectedFiles)

	// Stage selected but unstaged files
	if err = stageSelectedFiles(selectedFiles); err != nil {
		log.Error(err.Error())
		return
	}

	// Unstage the files that were not selected but are currently staged
	if err = unstageUnselectedFiles(unselectedFiles); err != nil {
		log.Error(err.Error())
		return
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
