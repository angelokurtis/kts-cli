package git

import (
	"strings"

	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func ListRemoteBranches() ([]string, error) {
	out, err := bash.RunAndLogRead("git branch -r --format '%(refname:short)'")
	if err != nil {
		return nil, err
	}

	branches := strings.Split(string(out), "\n")

	return branches, nil
}

func CurrentBranch() (string, error) {
	out, err := bash.RunAndLogRead("git branch --show-current")
	if err != nil {
		return "", err
	}

	branches := strings.Split(string(out), "\n")
	for _, branch := range branches {
		return branch, nil
	}

	return "", nil
}
