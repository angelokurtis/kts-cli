package format

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func createTemporaryDirectory(ctx context.Context) (string, func(), error) {
	slog.DebugContext(ctx, "Creating temporary directory")

	temp, err := os.MkdirTemp("", "kts-cli-format-")
	if err != nil {
		return "", func() {}, fmt.Errorf("failed to create temp directory: %w", err)
	}

	slog.InfoContext(ctx, "Temporary directory created", slog.String("tempDir", temp))

	return temp, func() {
		if err = os.RemoveAll(temp); err != nil {
			slog.WarnContext(ctx, "Failed to remove temporary directory", slog.String("tempDir", temp), tint.Err(err))
		} else {
			slog.InfoContext(ctx, "Temporary directory removed", slog.String("tempDir", temp))
		}
	}, nil
}
