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

func embedArgs(ctx context.Context, args ...any) (trace.Span, []any) {
	span := trace.SpanFromContext(ctx)
	sc := span.SpanContext()
	if sc.IsValid() {
		args = append(args, "trace_id", sc.TraceID().String(), "span_id", sc.SpanID().String())
		return span, args
	}
	return nil, args
}
