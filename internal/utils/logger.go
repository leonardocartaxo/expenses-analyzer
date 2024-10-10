package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"net/http"
	"os"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b) // Capture the response body
	return w.ResponseWriter.Write(b)
}

// LogRequestResponseMiddleware logs both the request and response
func LogRequestResponseMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract request headers
		headers := make(map[string][]string)
		for k, v := range c.Request.Header {
			headers[k] = v
		}

		// Read and log the request body (for POST and PUT requests)
		var requestBody interface{}
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				logger.Error("Error reading request body", "error", err)
				c.Next()
				return
			}

			// Try to parse the body as JSON
			if json.Unmarshal(bodyBytes, &requestBody) != nil {
				requestBody = string(bodyBytes) // If not valid JSON, log as string
			}

			// Reassign the body back to the request so it can be read again
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Capture the response body
		responseWriter := &responseBodyWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = responseWriter

		// Continue with the request processing
		c.Next()

		// Try to parse the response body as JSON
		var responseBody interface{}
		responseBytes := responseWriter.body.Bytes()
		if json.Unmarshal(responseBytes, &responseBody) != nil {
			responseBody = string(responseBytes) // If not valid JSON, log as string
		}

		// Log the request and response in structured and nested fields
		logger.InfoContext(context.TODO(), "Request and response",
			"request", slog.GroupValue(
				slog.String("method", c.Request.Method),
				slog.String("path", c.Request.URL.Path),
				slog.Any("headers", headers),
				slog.Any("body", requestBody),
			),
			"response", slog.GroupValue(
				slog.Int("status", c.Writer.Status()),
				slog.Any("body", responseBody),
			),
		)
	}
}

func NewLogger(logLevel int) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(logLevel)}))
}
