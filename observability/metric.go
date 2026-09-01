package observability

import (
	"context"
	"time"

	"github.com/BogdanBeliy/lgtm-obs/internal/otlp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type Metrics struct {
	Config   *Configs
	Resource *resource.Resource
	Provider *sdkmetric.MeterProvider

	exporter sdkmetric.Exporter
}

func (o *Observability) NewMetrics(ctx context.Context) (*Metrics, error) {
	exporter, err := otlp.NewMetricExporter(ctx, otlp.Endpoint{
		Protocol: o.Cfgs.Protocol,
		Path:     o.Cfgs.Path,
		Insecure: o.Cfgs.Insecure,
	})
	if err != nil {
		return nil, err
	}

	reader := sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(5*time.Second))

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(o.resource),
		sdkmetric.WithReader(reader))

	otel.SetMeterProvider(provider)

	return &Metrics{
		Config:   o.Cfgs,
		Resource: o.resource,
		Provider: provider,
		exporter: exporter,
	}, nil
}

func (m *Metrics) Shutdown(ctx context.Context) error {
	if m == nil || m.Provider == nil {
		return nil
	}
	return m.Provider.Shutdown(ctx)
}
