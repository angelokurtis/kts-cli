package git

import (
	"bufio"
	"bytes"
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

func ListStagedFiles() ([]string, error) {
	out, err := bash.RunAndLogRead("git diff --name-only --cached")
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		text := scanner.Text()
		files = append(files, text)
	}

	return files, err
}
