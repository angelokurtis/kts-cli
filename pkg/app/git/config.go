package git

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"github.com/angelokurtis/kts-cli/pkg/app/gpg"
	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func ConfigureSecretKey(sk *gpg.SecretKey) error {
	uid, ok := lo.First(sk.Uids)
	if !ok {
		return errors.New("")
	}

	steps := []string{
		fmt.Sprintf("git config user.name '%s'", uid.Name),
		fmt.Sprintf("git config user.email '%s'", uid.Email),
		"git config user.signingKey " + sk.KeyID,
		"git config commit.gpgsign true",
	}

	for _, cmd := range steps {
		if _, err := bash.RunAndLogWrite(cmd); err != nil {
			return err
		}
	}

	return nil
}
