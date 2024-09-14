package middlewares

import (
	"TradeEngine/API"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from the Authorization header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not provided"})
			c.Abort()
			return
		}

		// Extract userId using the external token microservice
		userId, err := API.ExtractUserIdFromToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token or failed to extract userId"})
			c.Abort()
			return
		}

		// Store userId in the context
		c.Set("userId", userId)

		// Proceed to the next middleware/handler
		c.Next()
	}
}
