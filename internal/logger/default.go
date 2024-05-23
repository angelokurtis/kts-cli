package logger

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"

	"github.com/angelokurtis/kts-cli/internal/otel"
)

const ktsLoglevelEnv = "KTS_LOGLEVEL"

// SetUp initializes the logger with predefined settings and returns the logger instance.
func SetUp() *slog.Logger {
	var lvl slog.Leveler = slog.LevelInfo

	// Adjust log level based on the environment variable value
	env, ok := os.LookupEnv(ktsLoglevelEnv)
	if ok {
		switch strings.ToLower(env) {
		case "debug":
			lvl = slog.LevelDebug
		case "info":
			lvl = slog.LevelInfo
		case "warn":
			lvl = slog.LevelWarn
		case "error":
			lvl = slog.LevelError
		}
	}

	// Create a new handler with tint colorized output
	handler := tint.NewHandler(os.Stderr, &tint.Options{
		AddSource:  true,
		Level:      lvl,
		TimeFormat: time.Kitchen,
	})

	// Wrap the tint handler with otel handler if otel integration is enabled (optional)
	handler = otel.NewLogHandler(handler)

	// Create a new slog logger instance with the formatted handler
	l := slog.New(handler)

	// Set the newly created logger as the default logger for the application
	slog.SetDefault(l)

	// Return the logger instance for potential further configuration
	return l
}
