// Package otlp creates OTLP exporters for observability signals.
package otlp

// Endpoint configures an OTLP exporter connection.
type Endpoint struct {
	Protocol string
	Path     string
	Insecure bool
}
