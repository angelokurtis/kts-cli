package git

import (
	"fmt"
	"strings"

	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func Stage(files []string) error {
	if len(files) == 0 {
		return nil
	}

	if _, err := bash.RunAndLogWrite(fmt.Sprintf("git add %s", strings.Join(files, " "))); err != nil {
		return err
	}

	return nil
}

func Unstage(files []string) error {
	if len(files) == 0 {
		return nil
	}

	if _, err := bash.RunAndLogWrite(fmt.Sprintf("git restore --staged %s", strings.Join(files, " "))); err != nil {
		return err
	}

	return nil
}
