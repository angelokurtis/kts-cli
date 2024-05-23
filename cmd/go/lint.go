package golang

import (
	log "log/slog"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/git"
	"github.com/angelokurtis/kts-cli/pkg/app/golangci_lint"
)

func lint(cmd *cobra.Command, args []string) {
	result, err := golangci_lint.Run()
	if err != nil {
		log.Error(err.Error())
		return
	}

	files, err := git.ShowDiffFiles("origin/HEAD")
	if err != nil {
		log.Error(err.Error())
		return
	}

	result, err = result.FilterByFiles(files)
	if err != nil {
		log.Error(err.Error())
		return
	}

	result.PrettyPrint()
}
