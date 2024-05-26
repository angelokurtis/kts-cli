package format

import (
	"context"
	"log/slog"
	"time"

	"github.com/lmittmann/tint"

	"github.com/angelokurtis/kts-cli/pkg/app/golang"
)

func runImportsReviser(ctx context.Context, workingDir string, fileArgs ...string) {
	runArgs := append([]string{"-set-alias", "-use-cache", "-rm-unused", "-format"}, fileArgs...)

	slog.DebugContext(ctx, "Running goimports-reviser")

	start := time.Now()
	if err := golang.RunSilently(workingDir, "github.com/incu6us/goimports-reviser/v3@latest", runArgs...); err != nil {
		elapsed := time.Since(start)
		slog.WarnContext(ctx, "Failed to run goimports-reviser", tint.Err(err), slog.Duration("duration", elapsed))
	} else {
		elapsed := time.Since(start)
		slog.InfoContext(ctx, "Successfully ran goimports-reviser", slog.Duration("duration", elapsed))
	}
}

func runGofumpt(ctx context.Context, workingDir string, fileArgs ...string) {
	runArgs := append([]string{"-w", "-extra"}, fileArgs...)

	slog.DebugContext(ctx, "Running gofumpt")

	start := time.Now()
	if err := golang.RunSilently(workingDir, "mvdan.cc/gofumpt@latest", runArgs...); err != nil {
		elapsed := time.Since(start)
		slog.WarnContext(ctx, "Failed to run gofumpt", tint.Err(err), slog.Duration("duration", elapsed))
	} else {
		elapsed := time.Since(start)
		slog.InfoContext(ctx, "Successfully ran gofumpt", slog.Duration("duration", elapsed))
	}
}

func runWsl(ctx context.Context, workingDir string, fileArgs ...string) {
	runArgs := append([]string{"-fix"}, fileArgs...)

	slog.DebugContext(ctx, "Running wsl")

	start := time.Now()
	if err := golang.RunSilently(workingDir, "github.com/bombsimon/wsl/v4/cmd...@latest", runArgs...); err != nil {
		elapsed := time.Since(start)
		slog.WarnContext(ctx, "Failed to run wsl", tint.Err(err), slog.Duration("duration", elapsed))
	} else {
		elapsed := time.Since(start)
		slog.InfoContext(ctx, "Successfully ran wsl", slog.Duration("duration", elapsed))
	}
}
