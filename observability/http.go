package observability

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func (o *Observability) WrapHandlers(handler http.Handler, operation string) http.Handler {
	return otelhttp.NewHandler(handler, operation)
}
