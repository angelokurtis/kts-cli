package golang

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	intlformat "github.com/angelokurtis/kts-cli/cmd/go/format"
	"github.com/angelokurtis/kts-cli/pkg/app/golang"
)

func clean(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current directory: %w", err)
	}

	pkgs, err := golang.ListPackages(wd)
	if err != nil {
		return fmt.Errorf("listing packages: %w", err)
	}

	source, err := intlformat.NewSourceCodes(wd, pkgs)
	if err != nil {
		return fmt.Errorf("loading source code: %w", err)
	}

	files, err := source.SelectMany()
	if err != nil {
		return err
	}

	for _, path := range files.RelativeFilePaths() {
		if err := removeComments(wd, path); err != nil {
			return err
		}

		slog.InfoContext(ctx, "Removed comments", slog.String("file", path))
	}

	return nil
}

func removeComments(wd, path string) error {
	// Define the shell script as a string
	shellScript := `
	#!/bin/bash
	
	set -e
	
	# Get the directory from the argument
	dir="$1"
	
	# Recursively find all .go files in the specified directory
	find "$dir" -type f -name "*.go" | while read -r file; do
		# Create a temp file
		tmp_file="${file}.tmp"
	
		# Remove comments and empty lines
		# - Reads whole file (-0777)
		# - Removes block comments /* ... */ safely
		# - Removes // comments only when outside of string literals
		# - Removes blank lines
		perl -0777 -pe '
			# --- Remove /* ... */ comments ---
			s{/\*.*?\*/}{}gs;
	
			# --- Remove // comments safely ---
			# This regex walks through each line, skipping // if it’s inside a quoted string.
			s{
				("(?:\\.|[^"\\])*") |  # Capture strings like "http://..." so we skip them
				(//[^\n]*)             # Capture actual line comments
			}{
				defined $1 ? $1 : ""   # Keep strings, remove comments
			}egmx;
	
			# --- Remove blank lines ---
			s/^\s*\n//mg;
		' "$file" > "$tmp_file"
	
		# Replace original file
		mv "$tmp_file" "$file"
	done
	`

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript, path)

	// Capture the output and error
	var stderr bytes.Buffer

	cmd.Dir = wd
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.MultiWriter(os.Stdout, &stderr)

	// Run the command
	if err := cmd.Run(); err != nil {
		return errors.Errorf("failed to remove file/dir: %s", strings.TrimSpace(stderr.String()))
	}

	return nil
}
