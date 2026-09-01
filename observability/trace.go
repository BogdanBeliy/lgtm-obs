package observability

import (
	"context"
	"time"

	"github.com/BogdanBeliy/lgtm-obs/internal/otlp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Traces struct {
	Config   *Configs
	Resource *resource.Resource
	Provider *sdktrace.TracerProvider

	exporter sdktrace.SpanExporter
}

func (o *Observability) NewTracing(ctx context.Context) (*Traces, error) {
	exporter, err := otlp.NewTraceExporter(ctx, otlp.Endpoint{
		Protocol: o.Cfgs.Protocol,
		Path:     o.Cfgs.Path,
		Insecure: o.Cfgs.Insecure,
	})
	if err != nil {
		return nil, err
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter,
			sdktrace.WithBatchTimeout(5*time.Second),
			sdktrace.WithMaxExportBatchSize(100)),
		sdktrace.WithResource(o.resource))

	otel.SetTracerProvider(provider)

	return &Traces{
		Config:   o.Cfgs,
		Resource: o.resource,
		Provider: provider,
		exporter: exporter,
	}, nil
}

func (t *Traces) Shutdown(ctx context.Context) error {
	if t == nil || t.Provider == nil {
		return nil
	}
	return t.Provider.Shutdown(ctx)
}
