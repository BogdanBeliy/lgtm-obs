// Example HTTP server with net/http instrumentation.
// Run: go run ./examples/http
package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/BogdanBeliy/lgtm-obs/observability"
	"go.opentelemetry.io/otel"
)

func main() {
	ctx := context.Background()

	obs, err := observability.InitObservability(ctx, &observability.Configs{
		Protocol:     "http",
		Path:         "localhost:4318",
		Insecure:     true,
		ExcludePaths: []string{"/health"},
		ServiceName:  "library-test",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer obs.Shutdown(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer("example").Start(r.Context(), "handle-hello")
		defer span.End()

		observability.Info(ctx, errors.New("boom").Error())
		_, _ = w.Write([]byte("Hello World"))
	})

	log.Println("http example on :8081")
	log.Fatal(http.ListenAndServe(":8081", obs.WrapHandlers(mux, "example-http")))
}
