package observability

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// NewLogger installs a slog logger that writes to stdout and OpenTelemetry.
func (o *Observability) NewLogger(logs *Logs) {
	otelHandler := otelslog.NewHandler("otel-slog", otelslog.WithLoggerProvider(logs.Provider))

	stdoutHandler := newStdoutHandler(o.Cfgs.LogFormat)

	mh := slog.NewMultiHandler(otelHandler, stdoutHandler)

	slog.SetDefault(slog.New(mh))
}

func newStdoutHandler(format string) slog.Handler {
	switch format {
	case "text":
		return slog.NewTextHandler(os.Stdout, nil)
	default:
		return slog.NewJSONHandler(os.Stdout, nil)
	}
}

// Info writes an informational log record and adds trace correlation fields.
func Info(ctx context.Context, msg string, args ...any) {
	_, a := embedArgs(ctx, args...)
	slog.InfoContext(ctx, msg, a...)
}

// Warn writes a warning log record and adds trace correlation fields.
func Warn(ctx context.Context, msg string, args ...any) {
	_, a := embedArgs(ctx, args...)
	slog.WarnContext(ctx, msg, a...)
}

// Error writes an error log record and marks the active span as failed.
func Error(ctx context.Context, msg string, args ...any) {
	span, a := embedArgs(ctx, args...)
	if span != nil {
		span.RecordError(errors.New(msg))
		span.SetStatus(codes.Error, msg)
	}
	slog.ErrorContext(ctx, msg, a...)
}

// AccessLogLevel maps HTTP status codes to slog levels:
// 200-399 INFO, 400-499 WARN, 500+ ERROR.
func AccessLogLevel(status int) slog.Level {
	switch {
	case status >= 500:
		return slog.LevelError
	case status >= 400:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

// LogAccess writes an access log with a unified message and level derived from status.
// accessErr is recorded on the span for 5xx responses when present.
func LogAccess(ctx context.Context, status int, msg string, accessErr error, args ...any) {
	_, a := embedArgs(ctx, args...)

	switch AccessLogLevel(status) {
	case slog.LevelInfo:
		slog.InfoContext(ctx, msg, a...)
	case slog.LevelWarn:
		slog.WarnContext(ctx, msg, a...)
	default:
		if span := trace.SpanFromContext(ctx); span != nil {
			if accessErr != nil {
				span.RecordError(accessErr)
				span.SetStatus(codes.Error, accessErr.Error())
			} else {
				span.SetStatus(codes.Error, msg)
			}
		}
		slog.ErrorContext(ctx, msg, a...)
	}
}

func embedArgs(ctx context.Context, args ...any) (trace.Span, []any) {
	span := trace.SpanFromContext(ctx)
	sc := span.SpanContext()
	if sc.IsValid() {
		args = append(args, "trace_id", sc.TraceID().String(), "span_id", sc.SpanID().String())
		return span, args
	}
	return nil, args
}
