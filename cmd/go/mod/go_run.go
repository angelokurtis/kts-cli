package mod

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func runGoModUpgrade(ctx context.Context, workingDir string) error {
	// Define the shell script as a string
	shellScript := `
	#!/bin/bash

	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	echo -e "${BLUE}go-mod-upgrade${NC}"
	go-mod-upgrade
	`

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var stderr bytes.Buffer

	cmd.Dir = workingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	// Run the command
	start := time.Now()

	if err := cmd.Run(); err != nil {
		// Check if the error is of type *exec.ExitError
		var exitError *exec.ExitError
		if errors.As(err, &exitError) && exitError.ExitCode() != 0 {
			return errors.New(strings.TrimSpace(stderr.String()))
		}
	}

	elapsed := time.Since(start)
	slog.DebugContext(ctx, "Successfully ran go-mod-upgrade", slog.Duration("duration", elapsed))

	return nil
}
