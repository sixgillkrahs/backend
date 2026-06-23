package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w responseBodyWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}

// ResponseWrapperMiddleware intercepts JSON responses and wraps them into {code, message, data} format.
func ResponseWrapperMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Do not wrap SSE (Server-Sent Events) streams
		if c.Writer.Header().Get("Content-Type") == "text/event-stream" || strings.HasPrefix(c.Request.URL.Path, "/api/v1/audit/stream") {
			c.Next()
			return
		}

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w
		c.Next()

		contentType := w.Header().Get("Content-Type")
		if strings.Contains(contentType, "application/json") && w.body.Len() > 0 {
			var jsonVal interface{}
			err := json.Unmarshal(w.body.Bytes(), &jsonVal)
			if err == nil {
				code := 0
				status := w.Status()
				if status >= 200 && status < 300 {
					code = 1
				}

				var message string
				var data interface{} = jsonVal

				// Extract standard errors or messages
				if m, ok := jsonVal.(map[string]interface{}); ok {
					if errMsg, exists := m["error"]; exists {
						message = fmt.Sprintf("%v", errMsg)
						data = nil
					} else if msg, exists := m["message"]; exists {
						message = fmt.Sprintf("%v", msg)
						// Keep message inside data or set it to null depending on user preference.
						// Usually message is a status message, so we can keep data as is or remove message.
					}
				}

				wrapped := gin.H{
					"code":    code,
					"message": message,
					"data":    data,
				}
				if code == 0 && message != "" {
					wrapped["error"] = message
				}

				newBody, err := json.Marshal(wrapped)
				if err == nil {
					w.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(len(newBody)))
					w.ResponseWriter.WriteHeader(status)
					w.ResponseWriter.Write(newBody)
					return
				}
			}
		}

		// Write original response if not JSON or parsing failed
		if w.body.Len() > 0 {
			w.ResponseWriter.Write(w.body.Bytes())
		}
	}
}
