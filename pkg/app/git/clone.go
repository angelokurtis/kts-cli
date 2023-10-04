package git

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/pkg/errors"
)

const suffix = ".git"

func Clone(repo string) error {
	dir, err := NewLocalDir(repo)
	if err != nil {
		return err
	}
	err = dir.CreateIfNotExist()
	if err != nil {
		return err
	}
	path := dir.Path()
	if dir.IsGithub() || dir.IsGitlab() {
		repo = dir.SSHAddress()
	}
	_, err = bash.RunAndLogWrite(fmt.Sprintf("git clone %s %s", repo, path))
	return errors.WithStack(err)
}
