package format

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func runImportsReviser(ctx context.Context, workingDir string, fileArgs ...string) error {
	// Define the shell script as a string
	args := strings.Join(fileArgs, " ")
	shellScript := fmt.Sprintf(`
	#!/bin/bash

	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	echo -e "${BLUE}goimports-reviser -set-alias -use-cache -rm-unused -format %s${NC}"
	goimports-reviser -set-alias -use-cache -rm-unused -format %s
	`, args, args)

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var stderr bytes.Buffer

	cmd.Dir = workingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	// Run the command
	start := time.Now()

	err := cmd.Run()
	if err != nil {
		// Check if the error is of type *exec.ExitError
		var exitError *exec.ExitError
		if errors.As(err, &exitError) && exitError.ExitCode() != 0 {
			return errors.New(strings.TrimSpace(stderr.String()))
		}
	}

	elapsed := time.Since(start)
	dirs := lo.Map(fileArgs, func(item string, index int) string { return filepath.Dir(item) })
	slog.DebugContext(ctx, "Successfully ran goimports-reviser", slog.Duration("duration", elapsed), slog.Any("paths", lo.Uniq(dirs)))

	return nil
}

func runGofumpt(ctx context.Context, workingDir string, fileArgs ...string) error {
	// Define the shell script as a string
	args := strings.Join(fileArgs, " ")
	shellScript := fmt.Sprintf(`
	#!/bin/bash

	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	echo -e "${BLUE}gofumpt -w -extra %s${NC}"
	gofumpt -w -extra %s
	`, args, args)

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var stderr bytes.Buffer

	cmd.Dir = workingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	start := time.Now()

	if err != nil {
		// Check if the error is of type *exec.ExitError
		var exitError *exec.ExitError
		if errors.As(err, &exitError) && exitError.ExitCode() != 0 {
			return errors.New(strings.TrimSpace(stderr.String()))
		}
	}

	elapsed := time.Since(start)
	dirs := lo.Map(fileArgs, func(item string, index int) string { return filepath.Dir(item) })
	slog.DebugContext(ctx, "Successfully ran gofumpt", slog.Duration("duration", elapsed), slog.Any("paths", lo.Uniq(dirs)))

	return nil
}

func runWsl(ctx context.Context, workingDir string, fileArgs ...string) error {
	// Define the shell script as a string
	args := strings.Join(fileArgs, " ")
	shellScript := fmt.Sprintf(`
	#!/bin/bash

	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	echo -e "${BLUE}wsl -force-err-cuddling -allow-cuddle-declarations -fix %s${NC}"
	wsl -force-err-cuddling -allow-cuddle-declarations -fix %s
	`, args, args)

	// Create a new command to run the script with the arguments
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var stderr bytes.Buffer

	cmd.Dir = workingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	// Run the command
	start := time.Now()

	err := cmd.Run()
	if err != nil {
		// Check if the error is of type *exec.ExitError
		var exitError *exec.ExitError
		if errors.As(err, &exitError) && exitError.ExitCode() != 0 && exitError.ExitCode() != 3 {
			return errors.New(strings.TrimSpace(stderr.String()))
		}
	}

	elapsed := time.Since(start)
	dirs := lo.Map(fileArgs, func(item string, index int) string { return filepath.Dir(item) })
	slog.DebugContext(ctx, "Successfully ran wsl", slog.Duration("duration", elapsed), slog.Any("paths", lo.Uniq(dirs)))

	return nil
}
