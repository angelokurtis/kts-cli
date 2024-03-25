package golang

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/git"
	"github.com/angelokurtis/kts-cli/pkg/app/golangci_lint"
)

func format(cmd *cobra.Command, args []string) {
	result, err := golangci_lint.Run()
	if err != nil {
		log.Fatal(err)
	}

	files, err := git.ShowDiffFiles("origin/HEAD")
	if err != nil {
		log.Fatal(err)
	}

	result, err = result.FilterByFiles(files)
	if err != nil {
		log.Fatal(err)
	}

	if err = result.PrettyPrint(); err != nil {
		log.Fatal(err)
	}
}
