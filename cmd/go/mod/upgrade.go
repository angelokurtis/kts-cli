package mod

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func upgrade(cmd *cobra.Command, args []string) error {
	// Create a new context
	ctx := context.Background()

	// Get the current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	if err = installGoModUpgrade(); err != nil {
		return err
	}

	if err = runGoModUpgrade(ctx, workingDir); err != nil {
		return err
	}

	return nil
}
