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

	if _, err := bash.RunAndLogWrite(fmt.Sprintf("git add -A -- %s", strings.Join(files, " "))); err != nil {
		return err
	}

	return nil
}

func Unstage(files []string) error {
	if len(files) == 0 {
		return nil
	}

	if _, err := bash.RunAndLogWrite(fmt.Sprintf("git restore --staged -- %s", strings.Join(files, " "))); err != nil {
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

func isUnstaged(status string) bool {
	if len(status) < 2 {
		return false // Incorrect or incomplete status input
	}

	// Get the second character which represents the working directory status
	workingDirectoryStatus := strings.TrimSpace(string(status[1]))

	// Check for any of the known unstaged statuses
	switch workingDirectoryStatus {
	case "M", "D", "?":
		return true
	default:
		return false
	}
}

func isStaged(status string) bool {
	if len(status) < 1 {
		return false // Incorrect status input
	}

	// Get the first character which represents the staging area status
	stagingStatus := strings.TrimSpace(string(status[0]))

	// Check for any of the known staged statuses
	switch stagingStatus {
	case "A", "M", "D", "R", "C":
		return true
	default:
		return false
	}
}
