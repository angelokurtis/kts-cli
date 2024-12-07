package files

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

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

	_, err := runListBySize(ctx, dir)
	check(err)
}

func runListBySize(ctx context.Context, workingDir string) ([]*file, error) {
	// Define the shell script as a string
	shellScript := `
	#!/bin/bash
	set -e
	du -ah --max-depth=1 | sort -h
	`

	color.Primary.Println("du -ah --max-depth=1 | sort -h")

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Dir = workingDir
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
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
	files := lo.FilterMap(lines, func(line string, index int) (*file, bool) {
		if len(line) < 1 {
			return nil, false
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) != 3 {
			return nil, false
		}

		return &file{raw: line, path: matches[2], size: matches[1]}, true
	})

	return files, nil
}

type file struct {
	raw  string
	path string
	size string
}
