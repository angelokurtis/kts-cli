package commits

import (
	"bytes"
	"context"
	log "log/slog"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var brazil *time.Location

func init() {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Error(err.Error())
		return
	}

	brazil = loc
}

// git commits list
func list(_ *cobra.Command, _ []string) {
	ctx := context.Background()

	err := runCommitList(ctx, dir)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Error(err.Error())
		return
	}
}

func runCommitList(ctx context.Context, workingDir string) error {
	// Define the shell script as a string
	shellScript := `
	#!/bin/bash

	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	echo -e "${BLUE}git log --graph --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit${NC}"
	git log --graph --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit
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
		return errors.Errorf("failed to list commits: %s", strings.TrimSpace(stderr.String()))
	}

	return nil
}
