package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"github.com/angelokurtis/kts-cli/pkg/bash"
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
	out, err := bash.RunAndLogRead("git status -s")
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		text := scanner.Text()
		splited := strings.Split(text, " ")
		splited = lo.Filter(splited, func(x string, index int) bool {
			trimed := strings.TrimSpace(x)
			return len(trimed) > 0
		})
		if len(splited) != 2 {
			continue
		}
		change := splited[0]
		file := splited[1]
		if change != "??" {
			files = append(files, file)
		}
	}

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

func CountCommitsByAuthor() (map[string]int64, error) {
	out, err := bash.RunAndLogRead("git shortlog --summary --email --numbered --all --no-merges")
	if err != nil {
		return nil, err
	}

	count := make(map[string]int64)
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.TrimSpace(text)
		splited := strings.Split(text, "\t")
		if len(splited) != 2 {
			continue
		}
		val, err := strconv.ParseInt(splited[0], 10, 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		count[splited[1]] = val
	}

	return count, nil
}
