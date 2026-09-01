package observability

import (
	"context"

	"github.com/BogdanBeliy/lgtm-obs/internal/otlp"
	g "go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/resource"
	sdklogs "go.opentelemetry.io/otel/sdk/log"
)

type Logs struct {
	Config   *Configs
	Resource *resource.Resource
	Provider *sdklogs.LoggerProvider

	exporter sdklogs.Exporter
}

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

func (l *Logs) Shutdown(ctx context.Context) error {
	if l == nil || l.Provider == nil {
		return nil
	}
	return l.Provider.Shutdown(ctx)
}
