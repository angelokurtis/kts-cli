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
	args := []string{
		"-set-alias",
		"-use-cache",
		"-rm-unused",
		"-format",
	}
	args = append(args, fileArgs...)
	argsStr := strings.Join(args, " ")

	shellScript := fmt.Sprintf(`
	#!/bin/bash

	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	echo -e "${BLUE}goimports-reviser %s${NC}"
	goimports-reviser %s
	`, argsStr, argsStr)

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
	args := []string{
		"-w",
		"-extra",
	}
	args = append(args, fileArgs...)
	argsStr := strings.Join(args, " ")

	shellScript := fmt.Sprintf(`
	#!/bin/bash

	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	echo -e "${BLUE}gofumpt %s${NC}"
	gofumpt %s
	`, argsStr, argsStr)

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
	enabled := []string{
		"after-block",
		"after-decl",
		"cuddle-group",
	}

	args := []string{
		"--default", "default",
		"--enable", strings.Join(enabled, ","),
		"--fix",
	}
	args = append(args, fileArgs...)
	argsStr := strings.Join(args, " ")

	shellScript := fmt.Sprintf(`
	#!/bin/bash

	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	echo -e "${BLUE}wsl %s${NC}"
	wsl %s
	`, argsStr, argsStr)

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

func runFuncorder(ctx context.Context, workingDir string, fileArgs ...string) error {
	// Define the shell script as a string
	args := []string{
		"-constructor=true",
		"-struct-method=true",
		"-alphabetical=true",
		"-function=true",
	}
	args = append(args, fileArgs...)
	argsStr := strings.Join(args, " ")

	shellScript := fmt.Sprintf(`
	#!/bin/bash

	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	echo -e "${BLUE}funcorder %s${NC}"
	funcorder %s
	`, argsStr, argsStr)

	// Create a new command to run the script
	cmd := exec.CommandContext(ctx, "bash", "-c", shellScript)

	// Capture the output and error
	var stderr bytes.Buffer

	cmd.Dir = workingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	// Run the command
	start := time.Now()

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, strings.TrimSpace(stderr.String()))
	}

	elapsed := time.Since(start)
	dirs := lo.Map(fileArgs, func(item string, index int) string { return filepath.Dir(item) })
	slog.DebugContext(ctx, "Successfully ran funcorder", slog.Duration("duration", elapsed), slog.Any("paths", lo.Uniq(dirs)))

	return nil
}

func runGoreorder(ctx context.Context, workingDir string, fileArgs ...string) error {
	// Define the shell script as a string
	args := []string{
		"reorder",
		"--write",
		"--reorder-types",
		"--format", "gofmt",
	}
	args = append(args, fileArgs...)
	argsStr := strings.Join(args, " ")

	shellScript := fmt.Sprintf(`
	#!/bin/bash

	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	echo -e "${BLUE}goreorder %s${NC}"
	goreorder %s
	`, argsStr, argsStr)

	// Create a new command to run the script
	cmd := exec.CommandContext(ctx, "bash", "-c", shellScript)

	// Capture the output and error
	var stderr bytes.Buffer

	cmd.Dir = workingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	// Run the command
	start := time.Now()

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, strings.TrimSpace(stderr.String()))
	}

	elapsed := time.Since(start)
	dirs := lo.Map(fileArgs, func(item string, index int) string { return filepath.Dir(item) })
	slog.DebugContext(ctx, "Successfully ran goreorder", slog.Duration("duration", elapsed), slog.Any("paths", lo.Uniq(dirs)))

	return nil
}
