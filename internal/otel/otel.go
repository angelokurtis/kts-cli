package otel

import (
	"context"
	"log/slog"
	"os"

	"github.com/angelokurtis/go-otel/starter"
	"github.com/lmittmann/tint"
)

// Environment variable names.
const (
	serviceNameEnv     = "OTEL_SERVICE_NAME"
	tracesExporterEnv  = "OTEL_TRACES_EXPORTER"
	metricsExporterEnv = "OTEL_METRICS_EXPORTER"
)

// SetUp sets up the OpenTelemetry instrumentation for the application.
// It returns the providers started by the starter package and a cleanup function
// to be called when the application exits.
func SetUp(ctx context.Context) (*starter.Providers, func()) {
	// If the environment variables for OpenTelemetry aren't set, set them to default values.
	if _, ok := os.LookupEnv(serviceNameEnv); !ok {
		_ = os.Setenv(serviceNameEnv, "kts-cli")
	}

	if _, ok := os.LookupEnv(tracesExporterEnv); !ok {
		_ = os.Setenv(tracesExporterEnv, "none")
	}

	if _, ok := os.LookupEnv(metricsExporterEnv); !ok {
		_ = os.Setenv(metricsExporterEnv, "none")
	}

	// Initiate the essential providers for OpenTelemetry.
	providers, shutdown, err := starter.StartProviders(ctx)
	if err != nil {
		slog.WarnContext(ctx, "Failed to start OpenTelemetry providers", tint.Err(err))
		return nil, func() {}
	}

	slog.DebugContext(ctx, "Started OpenTelemetry providers successfully")

	// Return the started providers and the cleanup function.
	return providers, func() {
		shutdown()
		slog.DebugContext(ctx, "Closed OpenTelemetry providers")
	}
}
