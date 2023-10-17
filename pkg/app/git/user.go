package git

import (
	"fmt"
	"strings"

	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func GetUser() (string, error) {
	name, err := GetUserName()
	if err != nil {
		return "", err
	}

	email, err := GetUserEmail()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s <%s>", name, email), nil
}

func GetUserName() (string, error) {
	out, err := bash.RunAndLogRead("git config user.name")
	if err != nil {
		return "", err
	}
	list := strings.Split(string(out), "\n")
	for _, name := range list {
		return name, nil
	}
	return "", nil
}

func GetUserEmail() (string, error) {
	out, err := bash.RunAndLogRead("git config user.email")
	if err != nil {
		return "", err
	}
	list := strings.Split(string(out), "\n")
	for _, email := range list {
		return email, nil
	}
	return "", nil
}
