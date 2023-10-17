package git

import (
	"os"
	"path/filepath"

	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/git"
)

func stage(cmd *cobra.Command, args []string) {
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

	staged, err := selectFiles(paths)
	if err != nil {
		log.Fatal(err)
	}

	if err = git.Stage(staged); err != nil {
		log.Fatal(err)
	}

	unstaged, _ := lo.Difference(paths, staged)

	if err = git.Unstage(unstaged); err != nil {
		log.Fatal(err)
	}
}
