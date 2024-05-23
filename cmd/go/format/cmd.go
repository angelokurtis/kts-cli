package format

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

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

	slog.DebugContext(ctx, "File exists and is not a directory", slog.String("filename", filename))

	if info.IsDir() {
		return fmt.Errorf("the specified path '%s' is a directory, expected a file", filename)
	}

	// Get absolute paths for the file and its directory
	absFilename, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to get absolute file path for '%s': %w", filename, err)
	}

	slog.DebugContext(ctx, "Obtained absolute file path", slog.String("absFilename", absFilename))

	absDir, err := filepath.Abs(path.Dir(filename))
	if err != nil {
		return fmt.Errorf("failed to get absolute path for directory '%s': %w", path.Dir(filename), err)
	}

	slog.DebugContext(ctx, "Obtained absolute directory path", slog.String("absDir", absDir))

	// Describe the Go package in the directory
	slog.DebugContext(ctx, "Describing Go package", slog.String("directory", absDir))

	pkgDetails, err := golang.DescribePackage(absDir)
	if err != nil {
		return fmt.Errorf("failed to describe package in directory '%s': %w", absDir, err)
	}

	slog.DebugContext(ctx, "Package described", slog.Int("GoFilesCount", len(pkgDetails.GoFiles)), slog.Int("TestGoFilesCount", len(pkgDetails.TestGoFiles)))
	files := append(pkgDetails.GoFiles, pkgDetails.TestGoFiles...)

	// Create a temporary directory and ensure it's removed afterwards
	slog.DebugContext(ctx, "Creating temporary directory")

	temp, err := os.MkdirTemp("", "kts-cli-format-")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	slog.InfoContext(ctx, "Temporary directory created", slog.String("tempDir", temp))

	defer func() {
		if err = os.RemoveAll(temp); err != nil {
			slog.WarnContext(ctx, "Failed to remove temporary directory", slog.String("tempDir", temp), tint.Err(err))
		} else {
			slog.InfoContext(ctx, "Temporary directory removed", slog.String("tempDir", temp))
		}
	}()

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Copy package files to the temporary directory
	slog.DebugContext(ctx, "Copying package files to temporary directory", slog.String("tempDir", temp))

	for _, file := range files {
		src := filepath.Join(absDir, file)
		dst := filepath.Join(temp, file)

		if err = copyFile(src, dst); err != nil {
			return fmt.Errorf("failed to copy file from '%s' to '%s': %w", src, dst, err)
		}

		srcRel, err := filepath.Rel(wd, src)
		if err != nil {
			return fmt.Errorf("failed to get relative path from '%s' to '%s': %w", wd, src, err)
		}

		slog.DebugContext(ctx, "Copied file", slog.String("src", "./"+srcRel), slog.String("dst", dst))
	}

	slog.DebugContext(ctx, "Obtained current working directory", slog.String("wd", wd))

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

	slog.DebugContext(ctx, "Prepared file arguments for formatting tools", slog.String("fileArgs", strings.Join(fileArgs, " ")))

	// Run goimports-reviser, gofumpt, and wsl on the files
	runArgs := append([]string{"-set-alias", "-use-cache", "-rm-unused", "-format"}, fileArgs...)

	slog.InfoContext(ctx, "Running goimports-reviser")

	if err = golang.RunSilently(wd, "github.com/incu6us/goimports-reviser/v3@latest", fileArgs...); err != nil {
		slog.WarnContext(ctx, "Failed to run goimports-reviser", tint.Err(err))
	} else {
		slog.InfoContext(ctx, "Successfully ran goimports-reviser")
	}

	runArgs = append([]string{"-w", "-extra"}, fileArgs...)

	slog.InfoContext(ctx, "Running gofumpt")

	if err = golang.RunSilently(wd, "mvdan.cc/gofumpt@latest", fileArgs...); err != nil {
		slog.WarnContext(ctx, "Failed to run gofumpt", tint.Err(err))
	} else {
		slog.InfoContext(ctx, "Successfully ran gofumpt")
	}

	runArgs = append([]string{"-fix"}, fileArgs...)

	slog.InfoContext(ctx, "Running wsl")

	if err = golang.RunSilently(wd, "github.com/bombsimon/wsl/v4/cmd...@latest", runArgs...); err != nil {
		slog.WarnContext(ctx, "Failed to run wsl", tint.Err(err))
	} else {
		slog.InfoContext(ctx, "Successfully ran wsl")
	}

	// Copy files back to the original directory
	slog.DebugContext(ctx, "Copying files back to original directory", slog.String("directory", absDir))

	for _, file := range files {
		src := filepath.Join(temp, file)

		dst := filepath.Join(absDir, file)
		if absFilename == dst {
			continue
		}

		if err = copyFile(src, dst); err != nil {
			return fmt.Errorf("failed to copy file from '%s' to '%s': %w", src, dst, err)
		}

		dstRel, err := filepath.Rel(wd, dst)
		if err != nil {
			return fmt.Errorf("failed to get relative path from '%s' to '%s': %w", wd, dst, err)
		}

		slog.InfoContext(ctx, "Copied file back", slog.String("src", src), slog.String("dst", "./"+dstRel))
	}

	slog.DebugContext(ctx, "Format run completed successfully", slog.String("filename", filename))

	return nil
}
