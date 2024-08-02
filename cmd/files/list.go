package files

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
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

func runListBySize(ctx context.Context, workingDir string) error {
	// Define the shell script as a string
	shellScript := `
	#!/bin/bash

	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	echo -e "${BLUE}du -hs * | sort -h${NC}"
	du -hs * | sort -h
	`

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var stderr bytes.Buffer

	cmd.Dir = workingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return errors.Errorf("failed to list files by size: %s", strings.TrimSpace(stderr.String()))
	}

	return nil
}
