package format

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

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
