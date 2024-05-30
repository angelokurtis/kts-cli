package format

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func runImportsReviser(ctx context.Context, workingDir string, fileArgs ...string) error {
	// Define the shell script as a string
	shellScript := fmt.Sprintf(`
	#!/bin/bash
	set -xe
	goimports-reviser -set-alias -use-cache -rm-unused -format %s
	`, strings.Join(fileArgs, " "))

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var out bytes.Buffer

	var stderr bytes.Buffer

	cmd.Stdout = &out
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
	slog.InfoContext(ctx, "Successfully ran goimports-reviser", slog.Duration("duration", elapsed))

	return nil
}

func runGofumpt(ctx context.Context, workingDir string, fileArgs ...string) error {
	// Define the shell script as a string
	shellScript := fmt.Sprintf(`
	#!/bin/bash
	set -xe
	gofumpt -w -extra %s
	`, strings.Join(fileArgs, " "))

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var out bytes.Buffer

	var stderr bytes.Buffer

	cmd.Stdout = &out
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
	slog.InfoContext(ctx, "Successfully ran gofumpt", slog.Duration("duration", elapsed))

	return nil
}

func runWsl(ctx context.Context, workingDir string, fileArgs ...string) error {
	// Define the shell script as a string
	shellScript := fmt.Sprintf(`
	#!/bin/bash
	set -xe
	wsl -force-err-cuddling -allow-cuddle-declarations -fix %s
	`, strings.Join(fileArgs, " "))

	// Create a new command to run the script with the arguments
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var out bytes.Buffer

	var stderr bytes.Buffer

	cmd.Stdout = &out
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
	slog.InfoContext(ctx, "Successfully ran wsl", slog.Duration("duration", elapsed))

	return nil
}
