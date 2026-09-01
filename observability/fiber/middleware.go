package fiber

import (
	"fmt"
	"net/http"
	"time"

	"github.com/BogdanBeliy/lgtm-obs/internal/paths"
	"github.com/BogdanBeliy/lgtm-obs/observability"
	"github.com/gofiber/fiber/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

type Config struct {
	ServiceName  string
	ExcludePaths []string
}

func ConfigFrom(cfg *observability.Configs) Config {
	if cfg == nil {
		return Config{}
	}
	return Config{
		ServiceName:  cfg.ServiceName,
		ExcludePaths: append([]string(nil), cfg.ExcludePaths...),
	}
}

func Middleware(cfg Config) fiber.Handler {
	if cfg.ServiceName == "" {
		return func(c fiber.Ctx) error { return c.Next() }
	}

	tr := otel.Tracer(cfg.ServiceName)

	return func(c fiber.Ctx) error {
		if paths.IsExcluded(c.Path(), cfg.ExcludePaths) {
			return c.Next()
		}

		started := time.Now()

		hdr := make(http.Header, len(c.GetReqHeaders()))
		for k, v := range c.GetReqHeaders() {
			for _, vv := range v {
				hdr.Add(k, vv)
			}
		}

		parentCtx := otel.GetTextMapPropagator().Extract(c.Context(), propagation.HeaderCarrier(hdr))
		spanName := fmt.Sprintf("%s %s", c.Method(), c.Path())
		ctx, span := tr.Start(parentCtx, spanName)
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.route", c.FullPath()),
			attribute.String("http.target", c.OriginalURL()),
			attribute.String("http.scheme", c.Protocol()),
		)

		c.SetContext(ctx)

		err := c.Next()

		status := c.Response().StatusCode()
		span.SetAttributes(attribute.Int("http.status_code", status))
		if status >= 500 {
			span.SetStatus(codes.Error, http.StatusText(status))
		}
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		logAccess(c, started, err)
		return err
	}
}
