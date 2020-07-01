package git

import (
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/angelokurtis/kts-cli/pkg/app/git"
	"github.com/angelokurtis/kts-cli/pkg/app/gpg"
	"github.com/spf13/cobra"
)

func signCommits(cmd *cobra.Command, args []string) {
	key, err := gpg.SelectSecretKey()
	if err != nil {
		common.Exit(err)
	}
	err = git.ConfigureSecretKey(key)
	if err != nil {
		common.Exit(err)
	}
}
