package git

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/git"
	"github.com/angelokurtis/kts-cli/pkg/app/gpg"
)

func signCommits(cmd *cobra.Command, args []string) {
	key, err := gpg.SelectSecretKey()
	if err != nil {
		system.Exit(err)
	}

	if key == nil {
		return
	}

	err = git.ConfigureSecretKey(key)
	if err != nil {
		system.Exit(err)
	}
}
