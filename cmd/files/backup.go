package files

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func backup(_ *cobra.Command, args []string) {
	ctx := context.Background()

	var dir string

	for i, arg := range args {
		switch i {
		case 0:
			dir = arg
		}
	}

	fileList, err := runGetBySize(ctx, dir)
	check(err)

	fileList = lo.Filter(fileList, func(item *file, index int) bool {
		return strings.TrimSpace(item.path) != "."
	})

	fileList, err = fileList.SelectFiles()
	check(err)

	for _, f := range fileList {
		info, err := os.Stat(f.path)
		check(err)

		if info.IsDir() {
			err = runS3Sync(ctx, dir, f.path)
			check(err)
		} else {
			err = runS3Cp(ctx, dir, f.path)
			check(err)
		}

		if rm {
			color.Primary.Println(fmt.Sprintf(`rm -rf "%s"`, f.path))
		}
	}
}

func runS3Cp(ctx context.Context, workingDir, path string) error {
	// Define the shell script as a string
	shellScript := fmt.Sprintf(`
	#!/bin/bash
	set -e
	aws s3 cp "%s" s3://kurtis-backups/
	`, path)

	color.Primary.Println(fmt.Sprintf(`aws s3 cp "%s" s3://kurtis-backups/`, path))

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var stderr bytes.Buffer

	cmd.Dir = workingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.MultiWriter(os.Stdout, &stderr)

	// Run the command
	if err := cmd.Run(); err != nil {
		return errors.Errorf("failed to remove file/dir: %s", strings.TrimSpace(stderr.String()))
	}

	return nil
}

func runS3Sync(ctx context.Context, workingDir, path string) error {
	// Define the shell script as a string
	shellScript := fmt.Sprintf(`
	#!/bin/bash
	set -e
	aws s3 sync "%s" s3://kurtis-backups/ --exact-timestamps
	`, path)

	color.Primary.Println(fmt.Sprintf(`aws s3 sync "%s" s3://kurtis-backups/ --exact-timestamps`, path))

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var stderr bytes.Buffer

	cmd.Dir = workingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.MultiWriter(os.Stdout, &stderr)

	// Run the command
	if err := cmd.Run(); err != nil {
		return errors.Errorf("failed to remove file/dir: %s", strings.TrimSpace(stderr.String()))
	}

	return nil
}
