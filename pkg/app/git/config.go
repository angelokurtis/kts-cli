package git

import (
	"github.com/angelokurtis/kts-cli/pkg/app/gpg"
	"strings"
)

func ConfigureSecretKey(sk *gpg.SecretKey) error {
	fields := strings.Fields(sk.UID)
	i := len(fields) - 1

	email := fields[i]
	email = strings.Replace(email, "<", "", -1)
	email = strings.Replace(email, ">", "", -1)

	name := strings.Join(fields[:i], " ")

	key := strings.Split(sk.Sec, "/")[1]

	_, err := runAndLog("config", "user.name", `"`+name+`"`)
	if err != nil {
		return err
	}
	_, err = runAndLog("config", "user.email", `"`+email+`"`)
	if err != nil {
		return err
	}
	_, err = runAndLog("config", "user.signingKey", key)
	if err != nil {
		return err
	}

	return nil
}
