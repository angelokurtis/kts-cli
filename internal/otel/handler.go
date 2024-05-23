package otel

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
)

// LogHandler wraps a slog.Handler and injects trace information into logs
type LogHandler struct {
	slog.Handler
}

// NewLogHandler creates a new LogHandler instance
func NewLogHandler(handler slog.Handler) *LogHandler {
	return &LogHandler{Handler: handler}
}

// Handle processes a log record and injects trace context if available
func (h *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	// Check if trace context exists in the context
	if span := trace.SpanContextFromContext(ctx); span.IsValid() {
		// Extract trace and span IDs from the context
		traceID := span.TraceID().String()
		spanID := span.SpanID().String()

		// Add trace and span IDs as attributes to the log record
		r.AddAttrs(
			slog.String("trace-id", traceID),
			slog.String("span-id", spanID),
		)
	}

	// Wrap the original handler's Handle call with error handling
	return errors.WithStack(h.Handler.Handle(ctx, r))
}
