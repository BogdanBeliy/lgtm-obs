package otlp

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// NewTraceExporter creates an OTLP trace exporter for the configured protocol.
func NewTraceExporter(ctx context.Context, cfg Endpoint) (sdktrace.SpanExporter, error) {
	switch cfg.Protocol {
	case "http":
		return newTraceHTTPExporter(ctx, cfg)
	case "grpc":
		return newTraceGRPCExporter(ctx, cfg)
	default:
		return nil, fmt.Errorf("observability: unsupported trace protocol %q", cfg.Protocol)
	}
}

func newTraceHTTPExporter(ctx context.Context, cfg Endpoint) (sdktrace.SpanExporter, error) {
	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(cfg.Path),
	}
	if cfg.Insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}
	return otlptracehttp.New(ctx, opts...)
}

func newTraceGRPCExporter(ctx context.Context, cfg Endpoint) (sdktrace.SpanExporter, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.Path),
	}
	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}
	return otlptracegrpc.New(ctx, opts...)
}
