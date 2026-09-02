// Example HTTP server with net/http instrumentation.
// Run: go run ./examples/http
package main

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/BogdanBeliy/lgtm-obs/observability"
)

func main() {
	ctx := context.Background()

	obs, err := observability.InitObservability(ctx, &observability.Configs{
		Protocol:     "http",
		Path:         "localhost:4318",
		Insecure:     true,
		ExcludePaths: []string{"/health"},
		ServiceName:  "library-test",
		SampleRate:   1.0,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer obs.Shutdown(ctx)

	client := &http.Client{
		Transport: obs.RoundTripper(nil),
		Timeout:   5 * time.Second,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/upstream", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("Hello from upstream"))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequestWithContext(
			r.Context(),
			http.MethodGet,
			"http://localhost:8081/upstream",
			nil,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			observability.Error(r.Context(), "upstream request failed", "error", err)
			http.Error(w, "upstream request failed", http.StatusBadGateway)
			return
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				slog.Error("Failed to close response body", "error", err)
			}
		}()
		_, _ = io.Copy(io.Discard, resp.Body)

		observability.Info(r.Context(), "upstream request completed", "status_code", resp.StatusCode)
		_, _ = w.Write([]byte("Hello World"))
	})

	log.Println("http example on :8081")
	server := &http.Server{
		Addr:              ":8081",
		Handler:           obs.WrapHandlers(mux, "example-http"),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
