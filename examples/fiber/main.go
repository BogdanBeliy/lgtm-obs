// Example Fiber server with otel middleware and HTTP access logs.
// Run: go run ./examples/fiber
package main

import (
	"context"
	"log"

	"github.com/BogdanBeliy/lgtm-obs/observability"
	obfiber "github.com/BogdanBeliy/lgtm-obs/observability/fiber"
	"github.com/gofiber/fiber/v3"
)

func main() {
	ctx := context.Background()

	obs, err := observability.InitObservability(ctx, &observability.Configs{
		Protocol:     "http",
		Path:         "localhost:4318",
		Insecure:     true,
		ExcludePaths: []string{"/health"},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer obs.Shutdown(ctx)

	app := fiber.New()
	app.Use(obfiber.Middleware(obfiber.ConfigFrom(obs.Cfgs)))

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello World")
	})
	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("ok")
	})

	log.Println("fiber example on :8083")
	log.Fatal(app.Listen(":8083"))
}
