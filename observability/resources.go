package observability

import (
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func (o *Observability) NewResource() *resource.Resource {
	res := resource.NewSchemaless(
		semconv.ServiceName(o.Cfgs.ServiceName),
		semconv.DeploymentEnvironment(o.Cfgs.Environment),
		semconv.ServiceVersion(o.Cfgs.Version))
	return res
}
