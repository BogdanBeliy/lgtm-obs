package observability

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/BogdanBeliy/lgtm-obs/internal/otlp"
)

// Metrics owns the OpenTelemetry meter provider and its configuration.
type Metrics struct {
	Config   *Configs
	Resource *resource.Resource
	Provider *sdkmetric.MeterProvider

	exporter sdkmetric.Exporter
}

// NewMetrics creates and installs an OTLP meter provider.
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

// Shutdown flushes pending metrics and stops the meter provider.
func (m *Metrics) Shutdown(ctx context.Context) error {
	if m == nil || m.Provider == nil {
		return nil
	}
	return m.Provider.Shutdown(ctx)
}
