// Example Gin server with otelgin middleware and HTTP access logs.
// Run: go run ./examples/gin
package main

import (
	"context"
	"log"

	"github.com/BogdanBeliy/lgtm-obs/observability"
	obgin "github.com/BogdanBeliy/lgtm-obs/observability/gin"
	"github.com/gin-gonic/gin"
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

	r := gin.New()
	r.Use(obgin.Middleware(obgin.ConfigFrom(obs.Cfgs)))

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	})

	log.Println("gin example on :8082")
	log.Fatal(r.Run(":8082"))
}
