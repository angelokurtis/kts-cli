package branch

import (
	"bytes"
	"context"
	slog "log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// git branch list
func list(_ *cobra.Command, _ []string) {
	ctx := context.Background()

	err := runBranchList(ctx, dir)
	check(err)
}

func check(err error) {
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func runBranchList(ctx context.Context, workingDir string) error {
	// Define the shell script as a string
	shellScript := `
	#!/bin/bash

	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	echo -e "${BLUE}git for-each-ref --sort=committerdate refs/heads/ --format='%(HEAD) %(color:yellow)%(refname:short)%(color:reset) - %(color:red)%(objectname:short)%(color:reset) - %(contents:subject) - %(authorname) (%(color:green)%(committerdate:relative)%(color:reset))'${NC}"
	git for-each-ref --sort=committerdate refs/heads/ --format='%(HEAD) %(color:yellow)%(refname:short)%(color:reset) - %(color:red)%(objectname:short)%(color:reset) - %(contents:subject) - %(authorname) (%(color:green)%(committerdate:relative)%(color:reset))'
	`

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var stderr bytes.Buffer

	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return errors.Errorf("failed to list branches: %s", strings.TrimSpace(stderr.String()))
	}

	return nil
}
