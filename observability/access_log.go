package observability

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

const maxResponseErrorBody = 256

// ResponseErrorMessage extracts a human-readable error from an HTTP response body.
func ResponseErrorMessage(status int, body []byte) string {
	if status < 400 || len(body) == 0 {
		return ""
	}

	body = bytes.TrimSpace(body)

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err == nil {
		for _, key := range []string{"error", "message", "detail", "title", "msg"} {
			if v, ok := payload[key]; ok {
				if s := fmt.Sprint(v); s != "" {
					return s
				}
			}
		}
	}

	s := string(body)
	if len(s) > maxResponseErrorBody {
		return s[:maxResponseErrorBody] + "..."
	}
	return s
}

// AccessLogError prefers a handler error and otherwise uses the response body for 4xx/5xx.
func AccessLogError(handlerErr error, status int, responseBody []byte) error {
	if handlerErr != nil {
		return handlerErr
	}
	if msg := ResponseErrorMessage(status, responseBody); msg != "" {
		return errors.New(msg)
	}
	return nil
}

// AccessLogMessage formats a one-line access log summary for log body display.
func AccessLogMessage(method, route, path string, status int, durationMs int64, err error) string {
	target := route
	if target == "" {
		target = path
	}
	msg := fmt.Sprintf("%s %s %d %dms", method, target, status, durationMs)
	if err != nil {
		return fmt.Sprintf("%s — %s", msg, err.Error())
	}
	return msg
}
