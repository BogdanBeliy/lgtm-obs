package observability

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/sdk/resource"
)

type Observability struct {
	Cfgs     *Configs
	Traces   *Traces
	Metrics  *Metrics
	Logs     *Logs
	resource *resource.Resource
}

func InitObservability(ctx context.Context, configs *Configs) (*Observability, error) {
	obs := &Observability{}

	obs.Cfgs = NewConfigs(configs)
	obs.resource = obs.NewResource()
	obs.NewPropagator()

	traces, err := obs.NewTracing(ctx)
	if err != nil {
		slog.Error("traces provider not running", "error", err)
		obs.Shutdown(ctx)
		return nil, err
	}
	obs.Traces = traces

	logs, err := obs.NewLogs(ctx)
	if err != nil {
		slog.Error("logs provider not running", "error", err)
		obs.Shutdown(ctx)
		return nil, err
	}
	obs.Logs = logs

	metrics, err := obs.NewMetrics(ctx)
	if err != nil {
		slog.Error("metrics provider not running", "error", err)
		obs.Shutdown(ctx)
		return nil, err
	}
	obs.Metrics = metrics

	obs.NewLogger(obs.Logs)
	return obs, nil
}

func (o *Observability) Shutdown(ctx context.Context) {
	if o.Traces != nil {
		if err := o.Traces.Shutdown(ctx); err != nil {
			slog.Error("traces provider shutdown error", "error", err)
		}
	}
	if o.Metrics != nil {
		if err := o.Metrics.Shutdown(ctx); err != nil {
			slog.Error("metrics provider shutdown error", "error", err)
		}
	}
	if o.Logs != nil {
		if err := o.Logs.Shutdown(ctx); err != nil {
			slog.Error("logs provider shutdown error", "error", err)
		}
	}
}
