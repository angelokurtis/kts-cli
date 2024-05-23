package format

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func Backup(ctx context.Context, srcFiles SourceCodes, tempDir string) error {
	if tempDir == "" {
		return fmt.Errorf("temporary directory path is empty")
	}

	info, err := os.Stat(tempDir)
	if err != nil {
		return fmt.Errorf("failed to stat temporary directory '%s': %w", tempDir, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("temporary path '%s' is not a directory", tempDir)
	}

	slog.DebugContext(ctx, "Copying package files to temporary directory", slog.String("directory", tempDir))

	for _, file := range srcFiles {
		src := file.FullFilePath()
		dst := filepath.Join(tempDir, file.RelativeFilePath())

		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("failed to copy file from '%s' to '%s': %w", src, dst, err)
		}

		slog.DebugContext(ctx, "Copied file",
			slog.String("src", "./"+file.RelativeFilePath()),
			slog.String("dst", dst),
		)
	}

	return nil
}

func Restore(ctx context.Context, backupFiles, selectedFiles SourceCodes, tempDir string) error {
	if tempDir == "" {
		return fmt.Errorf("temporary directory path is empty")
	}

	info, err := os.Stat(tempDir)
	if err != nil {
		return fmt.Errorf("failed to stat temporary directory '%s': %w", tempDir, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("temporary path '%s' is not a directory", tempDir)
	}

	slog.DebugContext(ctx, "Copying files back to original directory", slog.String("directory", tempDir))

	for _, file := range backupFiles {
		if selectedFiles.Contains(file) {
			continue
		}

		src := filepath.Join(tempDir, file.RelativeFilePath())
		dst := file.FullFilePath()

		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("failed to copy file from '%s' to '%s': %w", src, dst, err)
		}

		slog.DebugContext(ctx, "Copied file back",
			slog.String("src", src),
			slog.String("dst", "./"+file.RelativeFilePath()))
	}

	return nil
}
