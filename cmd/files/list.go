package files

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func list(_ *cobra.Command, args []string) {
	ctx := context.Background()

	var dir string

	for i, arg := range args {
		switch i {
		case 0:
			dir = arg
		}
	}

	err := runListBySize(ctx, dir)
	check(err)
}

func runRemove(ctx context.Context, workingDir, s string) error {
	// Define the shell script as a string
	shellScript := fmt.Sprintf(`
	#!/bin/bash
	set -e
	rm -rf %s
	`, s)

	color.Primary.Printf("rm -rf %s\n", s)

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var stderr bytes.Buffer

	cmd.Dir = workingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.MultiWriter(os.Stdout, &stderr)

	// Run the command
	if err := cmd.Run(); err != nil {
		return errors.Errorf("failed to remove file/dir: %s", strings.TrimSpace(stderr.String()))
	}

	return nil
}

func runListBySize(ctx context.Context, workingDir string) error {
	// Define the shell script as a string
	shellScript := `
	#!/bin/bash
	set -e
	du -ah --max-depth=1 | sort -h
	`

	color.Secondary.Println("du -ah --max-depth=1 | sort -h")

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var stderr bytes.Buffer

	cmd.Dir = workingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.MultiWriter(os.Stdout, &stderr)

	// Run the command
	if err := cmd.Run(); err != nil {
		return errors.Errorf("failed to list files by size: %s", strings.TrimSpace(stderr.String()))
	}

	return nil
}

func runGetBySize(ctx context.Context, workingDir string) (files, error) {
	// Define the shell script as a string
	shellScript := `
	#!/bin/bash
	set -e
	du -ah --max-depth=1 | sort -h
	`

	color.Secondary.Println("du -ah --max-depth=1 | sort -h")

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Dir = workingDir
	cmd.Stdout = &stdout
	cmd.Stderr = io.MultiWriter(os.Stdout, &stderr)

	// Run the command
	if err := cmd.Run(); err != nil {
		return nil, errors.Errorf("failed to list files by size: %s", strings.TrimSpace(stderr.String()))
	}

	return parseOutput(ctx, stdout)
}

func parseOutput(ctx context.Context, stdout bytes.Buffer) ([]*file, error) {
	scanner := bufio.NewScanner(&stdout)
	lines := make([]string, 0)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.Errorf("failed to scan output: %v", err)
	}

	re := regexp.MustCompile(`^(\S+)\s+(.+)$`)

	return lo.FilterMap(lines, func(line string, index int) (*file, bool) {
		if len(line) < 1 {
			return nil, false
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) != 3 {
			return nil, false
		}

		return &file{
			index: index,
			raw:   line,
			path:  matches[2],
			size:  matches[1],
		}, true
	}), nil
}

type files []*file

type file struct {
	index int
	raw   string
	path  string
	size  string
}

func (f files) SelectFiles() (files, error) {
	if len(f) == 0 {
		return f, nil
	}

	options := lo.KeyBy(f, func(item *file) string {
		return item.raw
	})
	keys := lo.Keys(options)
	sort.Slice(keys, func(i, j int) bool {
		opti := options[keys[i]]
		optj := options[keys[j]]

		return opti.index < optj.index
	})

	prompt := &survey.MultiSelect{
		Message: "Select the files to delete:",
		Options: keys,
	}

	var selects []string
	if err := survey.AskOne(prompt, &selects, survey.WithPageSize(10), survey.WithKeepFilter(true)); err != nil {
		return nil, errors.WithStack(err)
	}

	return lo.FilterMap(selects, func(item string, index int) (*file, bool) {
		selected, exists := options[item]
		return selected, exists
	}), nil
}
