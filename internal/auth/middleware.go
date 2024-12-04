package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AuthMiddleware checks for the Bearer token in the Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")

		// Check if the Authorization header contains the Bearer token
		if authHeader == "" || !strings.HasPrefix(authHeader, BearerSchema) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Extract the token from the header
		token := strings.TrimPrefix(authHeader, BearerSchema)

		// (Optional) Validate the token here - e.g., decode JWT or verify against database
		if !isTokenValid(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Proceed to the next handler if the token is valid
		c.Next()
	}
}

// Dummy token validation function
func isTokenValid(token string) bool {
	// In a real application, you would validate the token here
	// For example, check if it's a valid JWT or exists in a database
	return token == "valid-token" // Replace with actual token validation logic
}
