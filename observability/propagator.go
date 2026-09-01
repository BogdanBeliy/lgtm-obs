package observability

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// NewPropagator installs W3C Trace Context and Baggage propagation.
func (o *Observability) NewPropagator() {
	tc := propagation.TraceContext{}
	bg := propagation.Baggage{}
	ctmp := propagation.NewCompositeTextMapPropagator(tc, bg)
	otel.SetTextMapPropagator(ctmp)
}
