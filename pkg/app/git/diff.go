package git

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/pkg/errors"
	"os"
	"path"
	"path/filepath"
	"strconv"
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

func UncommittedFiles() ([]string, error) {
	out, err := bash.RunAndLogRead("git status -s | awk '{print $2}'")
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

	current, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(files))
	for _, file := range files {
		if len(file) == 0 {
			continue
		}
		result = append(result, filepath.Join(current, file))
	}

	return result, nil
}

func CountCommitsBetweenBranches(branch1, branch2 string) (int, error) {
	out, err := bash.RunAndLogRead(fmt.Sprintf("git rev-list --count %s ^%s", branch1, branch2))
	if err != nil {
		return 0, err
	}
	lines := strings.Split(string(out), "\n")
	var count int
	for _, line := range lines {
		count, err = strconv.Atoi(line)
		if err != nil {
			return 0, errors.WithStack(err)
		}
		break
	}
	return count, nil
}
