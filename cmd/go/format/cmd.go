package format

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/git"
	"github.com/angelokurtis/kts-cli/pkg/app/golang"
)

var filterUpdated = false
var all = false

var Command = &cobra.Command{
	Use: "format",
	Run: wrapWithErrorHandler(runFormat),
}

func init() {
	Command.PersistentFlags().BoolVarP(&filterUpdated, "updated", "u", false, "")
	Command.PersistentFlags().BoolVar(&all, "all", false, "")
}

func runFormat(cmd *cobra.Command, args []string) error {
	// Create a new context
	ctx := context.Background()

	// Get the current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	if all {
		runImportsReviser(ctx, workingDir, "./...")
		runGofumpt(ctx, workingDir, ".")
		runWsl(ctx, workingDir, "./...")

		return err
	}

	// List all Go packages in the current working directory, which will be processed
	pkgs, err := golang.ListPackages(workingDir)
	if err != nil {
		return fmt.Errorf("failed to list packages in the working directory: %w", err)
	}

	// Create a representation of the source code from the listed packages
	srcCodes, err := NewSourceCodes(workingDir, pkgs)
	if err != nil {
		return fmt.Errorf("failed to create source codes from the listed packages: %w", err)
	}

	selectedFiles := make(SourceCodes, 0)
	// Filter the source codes to include only uncommitted files if the --updated flag is enabled
	if filterUpdated {
		uncommittedFiles, err := git.UncommittedFiles()
		if err != nil {
			return fmt.Errorf("failed to retrieve uncommitted files: %w", err)
		}

		relativePaths, err := uncommittedFiles.RelativePaths()
		if err != nil {
			return fmt.Errorf("failed to get relative paths from uncommitted files: %w", err)
		}

		selectedFiles = srcCodes.FilterByRelativeFilePaths(relativePaths)
	} else {
		for _, srcCode := range srcCodes {
			selectedFiles = append(selectedFiles, srcCode)
		}
	}

	// Check if there are any source codes to process; return an error if none are found
	if len(selectedFiles) == 0 {
		return errors.New("no source codes found in the working directory")
	}

	// Select multiple source files for formatting
	selectedFiles, err = selectedFiles.SelectMany()
	if err != nil {
		return fmt.Errorf("failed to select multiple source files: %w", err)
	}

	selectedDirs := selectedFiles.RelativeDirPaths()

	affectedFiles := srcCodes.FilterByRelativeDirPaths(selectedDirs)

	// Create a temporary directory for backing up files before formatting and ensure that it is cleaned up after the operation
	tempDir, cleanup, err := createTemporaryDirectory(ctx)
	if err != nil {
		return fmt.Errorf("failed to create a temporary directory: %w", err)
	}
	defer cleanup()

	// Backup the affected files to the temporary directory
	if err = Backup(ctx, srcCodes, tempDir); err != nil {
		return fmt.Errorf("failed to backup affected files to the temporary directory: %w", err)
	}

	// Run goimports-reviser, gofumpt, and wsl on the files
	for _, relativeDirPath := range affectedFiles.RelativeDirPaths() {
		fileArgs := affectedFiles.FilterByRelativeDirPaths([]string{relativeDirPath}).RelativeFilePaths()
		runImportsReviser(ctx, workingDir, fileArgs...)
		runGofumpt(ctx, workingDir, fileArgs...)
		runWsl(ctx, workingDir, fileArgs...)
	}

	// Restore the formatted files back to their original locations
	if err = Restore(ctx, srcCodes, selectedFiles, tempDir); err != nil {
		return fmt.Errorf("failed to restore affected files from the temporary directory: %w", err)
	}

	slog.InfoContext(ctx, "Format run completed successfully")

	return nil
}
