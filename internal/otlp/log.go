package otlp

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	sdklogs "go.opentelemetry.io/otel/sdk/log"
)

func NewLogExporter(ctx context.Context, cfg Endpoint) (sdklogs.Exporter, error) {
	switch cfg.Protocol {
	case "http":
		return newLogHTTPExporter(ctx, cfg)
	case "grpc":
		return newLogGRPCExporter(ctx, cfg)
	default:
		return nil, fmt.Errorf("observability: unsupported logs protocol %q", cfg.Protocol)
	}
}

func newLogHTTPExporter(ctx context.Context, cfg Endpoint) (sdklogs.Exporter, error) {
	opts := []otlploghttp.Option{
		otlploghttp.WithEndpoint(cfg.Path),
	}
	if cfg.Insecure {
		opts = append(opts, otlploghttp.WithInsecure())
	}
	return otlploghttp.New(ctx, opts...)
}

func newLogGRPCExporter(ctx context.Context, cfg Endpoint) (sdklogs.Exporter, error) {
	opts := []otlploggrpc.Option{
		otlploggrpc.WithEndpoint(cfg.Path),
	}
	if cfg.Insecure {
		opts = append(opts, otlploggrpc.WithInsecure())
	}
	return otlploggrpc.New(ctx, opts...)
}
