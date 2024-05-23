package format

import (
	"context"
	"log/slog"

	"github.com/lmittmann/tint"

	"github.com/angelokurtis/kts-cli/pkg/app/golang"
)

func runImportsReviser(ctx context.Context, workingDir string, fileArgs ...string) {
	runArgs := append([]string{"-set-alias", "-use-cache", "-rm-unused", "-format"}, fileArgs...)

	slog.DebugContext(ctx, "Running goimports-reviser")

	if err := golang.RunSilently(workingDir, "github.com/incu6us/goimports-reviser/v3@latest", runArgs...); err != nil {
		slog.DebugContext(ctx, "Failed to run goimports-reviser", tint.Err(err))
	} else {
		slog.DebugContext(ctx, "Successfully ran goimports-reviser")
	}
}

func runGofumpt(ctx context.Context, workingDir string, fileArgs ...string) {
	runArgs := append([]string{"-w", "-extra"}, fileArgs...)

	slog.DebugContext(ctx, "Running gofumpt")

	if err := golang.RunSilently(workingDir, "mvdan.cc/gofumpt@latest", runArgs...); err != nil {
		slog.DebugContext(ctx, "Failed to run gofumpt", tint.Err(err))
	} else {
		slog.DebugContext(ctx, "Successfully ran gofumpt")
	}
}

func runWsl(ctx context.Context, workingDir string, fileArgs ...string) {
	runArgs := append([]string{"-fix"}, fileArgs...)

	slog.DebugContext(ctx, "Running wsl")

	if err := golang.RunSilently(workingDir, "github.com/bombsimon/wsl/v4/cmd...@latest", runArgs...); err != nil {
		slog.DebugContext(ctx, "Failed to run wsl", tint.Err(err))
	} else {
		slog.DebugContext(ctx, "Successfully ran wsl")
	}
}
