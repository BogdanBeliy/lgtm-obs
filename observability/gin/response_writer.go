package gin

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

const maxCapturedResponseBody = 4096

type responseBodyWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	if w.body.Len() < maxCapturedResponseBody {
		_, _ = w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w *responseBodyWriter) WriteString(s string) (int, error) {
	if w.body.Len() < maxCapturedResponseBody {
		_, _ = w.body.WriteString(s)
	}
	return w.ResponseWriter.WriteString(s)
}
