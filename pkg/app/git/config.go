package git

import (
	"strings"

	"github.com/angelokurtis/kts-cli/pkg/app/gpg"
	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func ConfigureSecretKey(sk *gpg.SecretKey) error {
	fields := strings.Fields(sk.UID)
	i := len(fields) - 1

	email := fields[i]
	email = strings.Replace(email, "<", "", -1)
	email = strings.Replace(email, ">", "", -1)

	name := strings.Join(fields[:i], " ")

	key := strings.Split(sk.Sec, "/")[1]

	if wordCount(name) > 1 {
		name = "'" + name + "'"
	}

	_, err := bash.RunAndLogWrite("git config user.name " + name)
	if err != nil {
		return err
	}

	_, err = bash.RunAndLogWrite("git config user.email " + email)
	if err != nil {
		return err
	}

	_, err = bash.RunAndLogWrite("git config user.signingKey " + key)
	if err != nil {
		return err
	}

	_, err = bash.RunAndLogWrite("git config commit.gpgsign true")
	if err != nil {
		return err
	}

	return nil
}

func wordCount(s string) int {
	words := strings.Fields(s)
	return len(words)
}
