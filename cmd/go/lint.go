package golang

import (
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/git"
	"github.com/angelokurtis/kts-cli/pkg/app/golangci_lint"
	"github.com/spf13/cobra"
)

func lint(cmd *cobra.Command, args []string) {
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

	result.PrettyPrint()
}
