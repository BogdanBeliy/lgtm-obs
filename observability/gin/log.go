package gin

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/BogdanBeliy/lgtm-obs/observability"
)

// RequestAttrs returns slog key-value pairs for an HTTP request (OpenTelemetry semconv names).
func RequestAttrs(c *gin.Context, started time.Time, accessErr error) []any {
	status := c.Writer.Status()
	attrs := []any{
		"http.method", c.Request.Method,
		"http.route", c.FullPath(),
		"http.target", c.Request.URL.String(),
		"http.scheme", c.Request.URL.Scheme,
		"http.status_code", status,
		"http.request_content_length", c.Request.ContentLength,
		"http.response_content_length", c.Writer.Size(),
		"client.address", c.ClientIP(),
		"user_agent.original", c.Request.UserAgent(),
		"http.request.header.referer", c.Request.Referer(),
		"url.path", c.Request.URL.Path,
		"url.query", c.Request.URL.RawQuery,
		"duration_ms", time.Since(started).Milliseconds(),
	}
	if accessErr != nil {
		attrs = append(attrs, "error.message", accessErr.Error())
	}
	return attrs
}

// InfoRequest logs an info message with HTTP request attributes.
func InfoRequest(c *gin.Context, started time.Time, msg string, args ...any) {
	args = append(RequestAttrs(c, started, nil), args...)
	observability.Info(c.Request.Context(), msg, args...)
}

// ErrorRequest logs an error with HTTP request attributes and records handler error.
func ErrorRequest(c *gin.Context, started time.Time, msg string, err error, args ...any) {
	args = append(RequestAttrs(c, started, err), args...)
	if msg == "" && err != nil {
		msg = err.Error()
	}
	observability.Error(c.Request.Context(), msg, args...)
}

func logAccess(c *gin.Context, started time.Time, handlerErr error, responseBody []byte) {
	status := c.Writer.Status()
	accessErr := observability.AccessLogError(handlerErr, status, responseBody)
	attrs := RequestAttrs(c, started, accessErr)
	msg := observability.AccessLogMessage(
		c.Request.Method,
		c.FullPath(),
		c.Request.URL.Path,
		status,
		time.Since(started).Milliseconds(),
		accessErr,
	)
	observability.LogAccess(c.Request.Context(), status, msg, accessErr, attrs...)
}
