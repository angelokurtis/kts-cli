package maxprocs

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"github.com/lmittmann/tint"
	"go.uber.org/automaxprocs/maxprocs"
)

// SetUp function to configure GOMAXPROCS and set up logging
func SetUp(logger *slog.Logger) (*SetupResult, func()) {
	ctx := context.Background()

	// Implementation for log output
	printf := func(format string, args ...any) {
		if !logger.Enabled(ctx, slog.LevelDebug) {
			return
		}

		var pcs [1]uintptr

		runtime.Callers(3, pcs[:])

		r := slog.NewRecord(time.Now(), slog.LevelDebug, fmt.Sprintf(format, args...), pcs[0])
		if err := logger.Handler().Handle(ctx, r); err != nil {
			slog.WarnContext(ctx, "Unable to handle log message. Please check the logger configuration and ensure that the necessary handler is set", tint.Err(err))
		}
	}

	// Setting up GOMAXPROCS and obtaining an undo function
	undo, err := maxprocs.Set(maxprocs.Logger(printf))
	if err != nil {
		slog.WarnContext(ctx, "Unable to automatically configure GOMAXPROCS", tint.Err(err))
	}

	// Returning the result of GOMAXPROCS configuration and the undo function
	return &SetupResult{successful: err != nil}, undo
}

// SetupResult holds the result of GOMAXPROCS configuration
type SetupResult struct{ successful bool }

// Successful checks if GOMAXPROCS configuration was successful
func (r *SetupResult) Successful() bool {
	return r.successful
}

// Failed checks if GOMAXPROCS configuration failed
func (r *SetupResult) Failed() bool {
	return !r.successful
}
