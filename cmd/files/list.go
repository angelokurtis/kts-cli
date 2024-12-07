package files

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"

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

func runListBySize(ctx context.Context, workingDir string) ([]string, error) {
	// Define the shell script as a string
	shellScript := `
	#!/bin/bash

	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	echo -e "${BLUE}du -ah --max-depth=1 | sort -h${NC}"
	du -ah --max-depth=1 | sort -h
	`

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

	return convertStdout(ctx, stdout)
}

func convertStdout(ctx context.Context, stdout bytes.Buffer) ([]string, error) {
	scanner := bufio.NewScanner(&stdout)
	res := make([]string, 0)

	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	lo.Filter(res, func(item string, index int) bool {
		return len(item) > 0
	})

	if err := scanner.Err(); err != nil {
		return nil, errors.Errorf("error occurred while scanning stdout: %v", err)
	}

	return nil, nil
}
