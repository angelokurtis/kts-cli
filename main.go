package main

import (
	"context"

	"github.com/angelokurtis/kts-cli/cmd"
	"github.com/angelokurtis/kts-cli/internal/logger"
	"github.com/angelokurtis/kts-cli/internal/maxprocs"
	"github.com/angelokurtis/kts-cli/internal/otel"
)

func main() {
	// Creating a context for the main function
	ctx := context.Background()

	// Set up logging
	l := logger.SetUp()

	// Set up GOMAXPROCS to utilize available CPU cores
	_, undo := maxprocs.SetUp(l)
	defer undo()

	// Set up tracing/metrics with OpenTelemetry
	_, shutdown := otel.SetUp(ctx)
	defer shutdown()

	cmd.Execute()
}
