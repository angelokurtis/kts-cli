package images

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func dedupe(cmd *cobra.Command, args []string) {
}

// TODO: Refactor this to use: 'rmlint /home/kurtis/Pictures/ --types=duplicates -S "Am" --output=json' instead.
// TODO: Consider moving this logic to 'kts files dedupe', 'kts files rm --duplicated', or possibly 'kts files dedupe rm'.
func runfindimagedupes(ctx context.Context, workingDir string) ([][]string, error) {
	// Define the shell script as a string
	shellScript := `
	#!/bin/bash
	set -e
	findimagedupes -q ./
	`

	color.Secondary.Printf("findimagedupes -q %s\n", workingDir)

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Dir = workingDir
	cmd.Stdout = &stdout
	cmd.Stderr = io.MultiWriter(os.Stdout, &stderr)

	// Run the command
	if err := cmd.Run(); err != nil {
		return nil, errors.Errorf(": %s", strings.TrimSpace(stderr.String()))
	}

	return parseOutput(ctx, stdout)
}

func parseOutput(ctx context.Context, buf bytes.Buffer) ([][]string, error) {
	var result [][]string

	// Read buffer as a string and split into lines
	lines := strings.Split(buf.String(), "\n")

	for _, line := range lines {
		// Trim spaces and skip empty lines
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Split the line by spaces to get file paths
		paths := strings.Fields(line)
		result = append(result, paths)
	}

	return result, nil
}
