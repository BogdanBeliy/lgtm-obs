package fiber

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/BogdanBeliy/lgtm-obs/observability"
)

// RequestAttrs returns slog key-value pairs for an HTTP request (OpenTelemetry semconv names).
func RequestAttrs(c fiber.Ctx, started time.Time, accessErr error) []any {
	status := c.Response().StatusCode()
	attrs := []any{
		"http.method", c.Method(),
		"http.route", c.FullPath(),
		"http.target", c.OriginalURL(),
		"http.scheme", c.Protocol(),
		"http.status_code", status,
		"http.request_content_length", requestContentLength(c),
		"http.response_content_length", len(c.Response().Body()),
		"client.address", c.IP(),
		"user_agent.original", c.Get("User-Agent"),
		"http.request.header.referer", c.Get("Referer"),
		"url.path", c.Path(),
		"url.query", string(c.Request().URI().QueryString()),
		"duration_ms", time.Since(started).Milliseconds(),
	}
	if accessErr != nil {
		attrs = append(attrs, "error.message", accessErr.Error())
	}
	return attrs
}

// InfoRequest logs an info message with HTTP request attributes.
func InfoRequest(c fiber.Ctx, started time.Time, msg string, args ...any) {
	args = append(RequestAttrs(c, started, nil), args...)
	observability.Info(c.Context(), msg, args...)
}

// ErrorRequest logs an error with HTTP request attributes and records handler error.
func ErrorRequest(c fiber.Ctx, started time.Time, msg string, err error, args ...any) {
	args = append(RequestAttrs(c, started, err), args...)
	if msg == "" && err != nil {
		msg = err.Error()
	}
	observability.Error(c.Context(), msg, args...)
}

func logAccess(c fiber.Ctx, started time.Time, handlerErr error) {
	status := c.Response().StatusCode()
	accessErr := observability.AccessLogError(handlerErr, status, c.Response().Body())
	attrs := RequestAttrs(c, started, accessErr)
	msg := observability.AccessLogMessage(
		c.Method(),
		c.FullPath(),
		c.Path(),
		status,
		time.Since(started).Milliseconds(),
		accessErr,
	)
	observability.LogAccess(c.Context(), status, msg, accessErr, attrs...)
}

func requestContentLength(c fiber.Ctx) int64 {
	raw := c.Get("Content-Length")
	if raw == "" {
		return 0
	}
	n, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0
	}
	return n
}
