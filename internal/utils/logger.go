package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
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

// SetRequestIDMiddleware to set a unique short request ID
func SetRequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, err := shortid.Generate()
		if err != nil {
			requestID = "unknown" // Fallback in case of error
		}
		c.Set("request_id", requestID)                   // Store request_id in the Gin context
		c.Writer.Header().Set("X-Request-ID", requestID) // Optionally add it to response headers
		c.Next()
	}
}

// LogRequestMiddleware logs the details of the request
func LogRequestMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract request headers
		headers := make(map[string][]string)
		for k, v := range c.Request.Header {
			headers[k] = v
		}

		// Read and parse the request body (for POST and PUT requests)
		var requestBody interface{}
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				logger.Error("Error reading request body", "error", err)
				return
			}

			// Try to parse the body as JSON
			if json.Unmarshal(bodyBytes, &requestBody) != nil {
				requestBody = string(bodyBytes) // Log as string if not valid JSON
			}

			// Reassign the body back to the request so it can be read again
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Log the request
		logger.InfoContext(context.TODO(), "Request received",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"headers", headers,
			"body", requestBody,
		)
	}
}

// LogResponseMiddleware logs the details of the response
func LogResponseMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capture the response body
		responseWriter := &responseBodyWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = responseWriter

		// Process the request
		c.Next()

		// Try to parse the response body as JSON
		var responseBody interface{}
		responseBytes := responseWriter.body.Bytes()
		if json.Unmarshal(responseBytes, &responseBody) != nil {
			responseBody = string(responseBytes) // Log as string if not valid JSON
		}

		// Log the response
		logger.InfoContext(context.TODO(), "Response sent",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"response_body", responseBody,
		)
	}
}

func NewLogger(logLevel int) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(logLevel)}))
}
