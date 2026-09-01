package gin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/BogdanBeliy/lgtm-obs/internal/paths"
	"github.com/BogdanBeliy/lgtm-obs/observability"
)

// Config defines Gin middleware settings.
type Config struct {
	ServiceName  string
	ExcludePaths []string
}

// ConfigFrom creates Gin middleware settings from observability configuration.
func ConfigFrom(cfg *observability.Configs) Config {
	if cfg == nil {
		return Config{}
	}
	return Config{
		ServiceName:  cfg.ServiceName,
		ExcludePaths: append([]string(nil), cfg.ExcludePaths...),
	}
}

// Middleware returns Gin middleware for tracing and access logging.
func Middleware(cfg Config) gin.HandlerFunc {
	if cfg.ServiceName == "" {
		return func(c *gin.Context) { c.Next() }
	}

	opts := []otelgin.Option{}
	if len(cfg.ExcludePaths) > 0 {
		opts = append(opts, otelgin.WithFilter(excludeFilter(cfg.ExcludePaths)))
	}

	otelMW := otelgin.Middleware(cfg.ServiceName, opts...)

	return func(c *gin.Context) {
		if paths.IsExcluded(c.Request.URL.Path, cfg.ExcludePaths) {
			c.Next()
			return
		}

		started := time.Now()
		bodyWriter := &responseBodyWriter{ResponseWriter: c.Writer}
		c.Writer = bodyWriter
		otelMW(c)

		var handlerErr error
		if last := c.Errors.Last(); last != nil {
			handlerErr = last.Err
		}
		logAccess(c, started, handlerErr, bodyWriter.body.Bytes())
	}
}

func excludeFilter(exclude []string) otelgin.Filter {
	return func(r *http.Request) bool {
		return !paths.IsExcluded(r.URL.Path, exclude)
	}
}
