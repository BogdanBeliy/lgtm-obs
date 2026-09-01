package observability

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

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

func Info(ctx context.Context, msg string, args ...any) {
	_, a := embedArgs(ctx, args...)
	slog.InfoContext(ctx, msg, a...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	_, a := embedArgs(ctx, args...)
	slog.WarnContext(ctx, msg, a...)
}

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
// handlerErr is recorded on the span for 5xx responses when present.
func LogAccess(ctx context.Context, status int, msg string, handlerErr error, args ...any) {
	_, a := embedArgs(ctx, args...)

	switch AccessLogLevel(status) {
	case slog.LevelInfo:
		slog.InfoContext(ctx, msg, a...)
	case slog.LevelWarn:
		slog.WarnContext(ctx, msg, a...)
	default:
		if span := trace.SpanFromContext(ctx); span != nil {
			if handlerErr != nil {
				span.RecordError(handlerErr)
				span.SetStatus(codes.Error, handlerErr.Error())
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

func AccessLogMessage(method, route, path string, status int, durationMs int64) string {
	target := route
	if target == "" {
		target = path
	}
	return fmt.Sprintf("%s %s %d %dms", method, target, status, durationMs)
}
