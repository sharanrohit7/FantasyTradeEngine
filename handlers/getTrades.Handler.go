package handlers

import (
	"TradeEngine/services"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// GetTradesHandler handles requests to retrieve trades for a given userId.
func GetTradesHandler(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
			return
		}
		userIdStr, ok := userId.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userId type"})
			return
		}
		// Retrieve trades from Redis
		trades := services.GetTradesByUserId(context.Background(), rdb, userIdStr)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 	return
		// }

		// Respond with trades data
		c.JSON(http.StatusOK, gin.H{"trades": trades})
	}
}
