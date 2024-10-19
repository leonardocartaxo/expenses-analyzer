package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
	"io"
	"log/slog"
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

// Middleware to set a unique short request ID
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
		requestID, _ := c.Get("request_id") // Retrieve the request_id from context

		// Check if the Content-Type is JSON
		if c.Request.Header.Get("Content-Type") == "application/json" {
			// Read and parse the request body as JSON
			var requestBody interface{}
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				logger.Error("Error reading request body", "request_id", requestID, "error", err)
				return
			}

			// Try to parse the body as JSON
			if json.Unmarshal(bodyBytes, &requestBody) == nil {
				// Log the request with the parsed JSON body
				logger.InfoContext(context.TODO(), "Request received",
					"request_id", requestID,
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"body", requestBody,
				)
			}

			// Reassign the body back to the request so it can be read again
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}
}

// LogResponseMiddleware logs the details of the response
func LogResponseMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, _ := c.Get("request_id") // Retrieve the request_id from context

		// Capture the response body
		responseWriter := &responseBodyWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = responseWriter

		// Process the request
		c.Next()

		// Check if the response Content-Type is JSON
		if c.Writer.Header().Get("Content-Type") == "application/json" {
			// Try to parse the response body as JSON
			var responseBody interface{}
			responseBytes := responseWriter.body.Bytes()
			if json.Unmarshal(responseBytes, &responseBody) == nil {
				// Log the response with the parsed JSON body
				logger.InfoContext(context.TODO(), "Response sent",
					"request_id", requestID,
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"status", c.Writer.Status(),
					"response_body", responseBody,
				)
			}
		}
	}
}

func NewLogger(logLevel int) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(logLevel)}))
}
