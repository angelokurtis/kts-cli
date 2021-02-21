package cmd

import (
	"github.com/angelokurtis/kts-cli/pkg/app/git"
	"github.com/angelokurtis/kts-cli/pkg/app/gpg"
	"github.com/spf13/cobra"
)

// gitSignCmd represents the git sign command
var gitSignCmd = &cobra.Command{
	Use:   "sign",
	Short: "Tells Git about your signing key to sign commits",
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := gpg.SelectSecretKey()
		if err != nil {
			return err
		}
		if key == nil {
			return nil
		}
		err = git.ConfigureSecretKey(key)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	gitCmd.AddCommand(gitSignCmd)
}
