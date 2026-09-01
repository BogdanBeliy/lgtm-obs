package otlp

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// NewMetricExporter creates an OTLP metric exporter for the configured protocol.
func NewMetricExporter(ctx context.Context, cfg Endpoint) (sdkmetric.Exporter, error) {
	switch cfg.Protocol {
	case "http":
		return newMetricHTTPExporter(ctx, cfg)
	case "grpc":
		return newMetricGRPCExporter(ctx, cfg)
	default:
		return nil, fmt.Errorf("observability: unsupported metrics protocol %q", cfg.Protocol)
	}
}

func newMetricHTTPExporter(ctx context.Context, cfg Endpoint) (sdkmetric.Exporter, error) {
	opts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(cfg.Path),
	}
	if cfg.Insecure {
		opts = append(opts, otlpmetrichttp.WithInsecure())
	}
	return otlpmetrichttp.New(ctx, opts...)
}

func newMetricGRPCExporter(ctx context.Context, cfg Endpoint) (sdkmetric.Exporter, error) {
	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(cfg.Path),
	}
	if cfg.Insecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}
	return otlpmetricgrpc.New(ctx, opts...)
}
