// middleware/token_verify.go
package middleware

import (
	"api-gateway/config"
	"bytes"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// TokenVerifyMiddleware verifies the JWT token by calling the Core Service
func TokenVerifyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Extract the token from the header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Call the Core Service to verify the token
		resp, err := http.Post(
			config.CoreServiceURL+"/auth/verify-token",
			"application/json",
			bytes.NewBuffer([]byte(`{"token":"`+tokenString+`"}`)),
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to communicate with Core Service"})
			c.Abort()
			return
		}

		defer resp.Body.Close()

		// Check the response from the Core Service
		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Token is valid, proceed to route the request
		c.Next()
	}
}
