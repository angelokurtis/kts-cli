package git

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"path"
	"strings"
)

func ShowDiffFiles(branch string) ([]string, error) {
	out, err := bash.RunAndLogRead(fmt.Sprintf("git diff --name-only %s", branch))
	if err != nil {
		return nil, err
	}

	files := strings.Split(string(out), "\n")

	if len(files) == 0 {
		return nil, nil
	}

	out, err = bash.Run("git rev-parse --show-toplevel")
	if err != nil {
		return nil, err
	}

	toplevels := strings.Split(string(out), "\n")
	var toplevel string
	for _, tl := range toplevels {
		toplevel = tl
		break
	}

	result := make([]string, 0, len(files))
	for _, file := range files {
		if len(file) == 0 {
			continue
		}
		result = append(result, path.Join(toplevel, file))
	}

	return result, nil
}
