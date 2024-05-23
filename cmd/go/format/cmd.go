package format

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/lmittmann/tint"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/golang"
)

var Command = &cobra.Command{
	Use: "format",
	Run: wrapWithErrorHandler(runFormat),
}

func runFormat(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	slog.InfoContext(ctx, "Running format command")

	var path string
	if len(args) > 0 {
		path = args[0]
	}

	if path == "" {
		return errors.New("no path provided")
	}

	slog.InfoContext(ctx, "Describing package", slog.String("path", path))

	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	slog.InfoContext(ctx, "Obtained absolute path", slog.String("absPath", abs))

	pkgDetails, err := golang.DescribePackage(abs)
	if err != nil {
		return fmt.Errorf("failed to describe package: %w", err)
	}

	slog.InfoContext(ctx, "Package described successfully", slog.String("path", path))

	files := append(pkgDetails.GoFiles, pkgDetails.TestGoFiles...)

	slog.InfoContext(ctx, "Creating temporary directory for format operation")

	temp, err := os.MkdirTemp("", "kts-cli-format-")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	defer func() {
		if err = os.RemoveAll(temp); err != nil {
			slog.WarnContext(ctx, "Failed to remove temporary directory", slog.String("tempDir", temp), tint.Err(err))
			return
		}

		slog.InfoContext(ctx, "Temporary directory removed successfully", slog.String("tempDir", temp))
	}()

	slog.InfoContext(ctx, "Temporary directory created", slog.String("tempDir", temp))

	for _, file := range files {
		src := filepath.Join(abs, file)
		dst := filepath.Join(temp, file)

		slog.InfoContext(ctx, "Copying file", slog.String("src", src), slog.String("dst", dst))

		if err = copyFile(src, dst); err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}

		slog.InfoContext(ctx, "File copied successfully", slog.String("src", src), slog.String("dst", dst))
	}

	slog.InfoContext(ctx, "Format command completed successfully")

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
