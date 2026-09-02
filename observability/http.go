package observability

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// WrapHandlers instruments an HTTP handler with OpenTelemetry.
func (o *Observability) WrapHandlers(handler http.Handler, operation string) http.Handler {
	return otelhttp.NewHandler(handler, operation)
}

// RoundTripper instruments outgoing HTTP requests with OpenTelemetry.
// A nil base transport uses http.DefaultTransport.
func (o *Observability) RoundTripper(base http.RoundTripper) http.RoundTripper {
	return otelhttp.NewTransport(base)
}
