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
	ctx := context.Background()

	var filename string
	if len(args) > 0 {
		filename = args[0]
	}

	info, err := os.Stat(filename)
	if err != nil {
		return fmt.Errorf("failed to retrieve file information for path '%s': %w", filename, err)
	}

	if info.IsDir() {
		return fmt.Errorf("the specified path '%s' is a directory, expected a file", filename)
	}

	absFilename, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to get absolute file path for '%s': %w", filename, err)
	}

	dir := path.Dir(filename)

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for directory '%s': %w", dir, err)
	}

	pkgDetails, err := golang.DescribePackage(absDir)
	if err != nil {
		return fmt.Errorf("failed to describe package in directory '%s': %w", absDir, err)
	}

	files := append(pkgDetails.GoFiles, pkgDetails.TestGoFiles...)

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

	for _, file := range files {
		src := filepath.Join(absDir, file)
		dst := filepath.Join(temp, file)

		if err = copyFile(src, dst); err != nil {
			return fmt.Errorf("failed to copy file from '%s' to '%s': %w", src, dst, err)
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	fileArgs := make([]string, 0, len(files))

	for _, file := range files {
		f := filepath.Join(absDir, file)

		rel, err := filepath.Rel(wd, f)
		if err != nil {
			return fmt.Errorf("failed to get relative path from '%s' to '%s': %w", wd, f, err)
		}

		fileArgs = append(fileArgs, rel)
	}

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
