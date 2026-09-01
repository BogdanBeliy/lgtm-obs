package observability

import (
	"context"

	g "go.opentelemetry.io/otel/log/global"
	sdklogs "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/BogdanBeliy/lgtm-obs/internal/otlp"
)

// Logs owns the OpenTelemetry logger provider and its configuration.
type Logs struct {
	Config   *Configs
	Resource *resource.Resource
	Provider *sdklogs.LoggerProvider

	exporter sdklogs.Exporter
}

// NewLogs creates and installs an OTLP logger provider.
func (o *Observability) NewLogs(ctx context.Context) (*Logs, error) {
	exporter, err := otlp.NewLogExporter(ctx, otlp.Endpoint{
		Protocol: o.Cfgs.Protocol,
		Path:     o.Cfgs.Path,
		Insecure: o.Cfgs.Insecure,
	})
	if err != nil {
		return nil, err
	}

	processor := sdklogs.NewBatchProcessor(exporter)

	provider := sdklogs.NewLoggerProvider(
		sdklogs.WithProcessor(processor),
		sdklogs.WithResource(o.resource))

	g.SetLoggerProvider(provider)

	return &Logs{
		Config:   o.Cfgs,
		Resource: o.resource,
		Provider: provider,
		exporter: exporter,
	}, nil
}

// Shutdown flushes pending log records and stops the logger provider.
func (l *Logs) Shutdown(ctx context.Context) error {
	if l == nil || l.Provider == nil {
		return nil
	}
	return l.Provider.Shutdown(ctx)
}
