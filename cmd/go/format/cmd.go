package format

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"

	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/golang"
)

var Command = &cobra.Command{
	Use: "format",
	Run: wrapWithErrorHandler(runFormat),
}

func runFormat(cmd *cobra.Command, args []string) error {
	// Create a new context
	ctx := context.Background()

	// Determine the filename from arguments
	var filename string
	if len(args) > 0 {
		filename = args[0]
	}

	// Check if the file exists and is not a directory
	info, err := os.Stat(filename)
	if err != nil {
		return fmt.Errorf("failed to retrieve file information for path '%s': %w", filename, err)
	}

	if info.IsDir() {
		return fmt.Errorf("the specified path '%s' is a directory, expected a file", filename)
	}

	// Get absolute paths for the file and its directory
	absFilename, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to get absolute file path for '%s': %w", filename, err)
	}

	absDir, err := filepath.Abs(path.Dir(filename))
	if err != nil {
		return fmt.Errorf("failed to get absolute path for directory '%s': %w", path.Dir(filename), err)
	}

	// Describe the Go package in the directory
	pkgDetails, err := golang.DescribePackage(absDir)
	if err != nil {
		return fmt.Errorf("failed to describe package in directory '%s': %w", absDir, err)
	}

	files := append(pkgDetails.GoFiles, pkgDetails.TestGoFiles...)

	// Create a temporary directory and ensure it's removed afterwards
	temp, err := os.MkdirTemp("", "kts-cli-format-")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	defer func() {
		if err = os.RemoveAll(temp); err != nil {
			slog.WarnContext(ctx, "Failed to remove temporary directory", slog.String("tempDir", temp), tint.Err(err))
			return
		}
	}()

	// Copy package files to the temporary directory
	for _, file := range files {
		src := filepath.Join(absDir, file)
		dst := filepath.Join(temp, file)

		if err = copyFile(src, dst); err != nil {
			return fmt.Errorf("failed to copy file from '%s' to '%s': %w", src, dst, err)
		}
	}

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Prepare file arguments for formatting tools
	fileArgs := make([]string, 0, len(files))
	for _, file := range files {
		f := filepath.Join(absDir, file)

		rel, err := filepath.Rel(wd, f)
		if err != nil {
			return fmt.Errorf("failed to get relative path from '%s' to '%s': %w", wd, f, err)
		}

		fileArgs = append(fileArgs, rel)
	}

	// Run goimports-reviser, gofumpt, and wsl on the files
	runArgs := append([]string{"-set-alias", "-use-cache", "-rm-unused", "-format"}, fileArgs...)

	if err = golang.Run(wd, "github.com/incu6us/goimports-reviser/v3@latest", fileArgs...); err != nil {
		slog.WarnContext(ctx, "Failed to run goimports-reviser")
	}

	runArgs = append([]string{"-w", "-extra"}, fileArgs...)

	if err = golang.Run(wd, "mvdan.cc/gofumpt@latest", fileArgs...); err != nil {
		slog.WarnContext(ctx, "Failed to run gofumpt")
	}

	runArgs = append([]string{"-fix"}, fileArgs...)
	if err = golang.Run(wd, "github.com/bombsimon/wsl/v4/cmd...@latest", runArgs...); err != nil {
		slog.WarnContext(ctx, "Failed to run wsl")
	}

	// Copy files back to the original directory
	for _, file := range files {
		src := filepath.Join(temp, file)

		dst := filepath.Join(absDir, file)
		if absFilename == dst {
			continue
		}

		if err = copyFile(src, dst); err != nil {
			return fmt.Errorf("failed to copy file from '%s' to '%s': %w", src, dst, err)
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}
