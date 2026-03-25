package actions

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type RequestLogger struct {
	DefaultTransport http.RoundTripper
}

func (rl *RequestLogger) RoundTrip(r *http.Request) (*http.Response, error) {
	path, err := url.PathUnescape(r.URL.Path)
	if err != nil {
		path = r.URL.Path
	}

	rawQuery, err := url.QueryUnescape(r.URL.RawQuery)
	if err != nil {
		rawQuery = r.URL.RawQuery
	}

	if rawQuery != "" {
		path += "?" + rawQuery
	}

	args := []any{
		"request", fmt.Sprintf("%s %s", r.Method, path),
	}
	startTime := time.Now()

	resp, err := rl.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	endTime := time.Now()
	args = append(args,
		"response", resp.Status,
		"duration", endTime.Sub(startTime),
	)
	ctx := r.Context()
	slog.InfoContext(ctx, "call to API", args...)

	return resp, nil
}
