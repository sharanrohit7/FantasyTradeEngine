package handlers

import (
	"TradeEngine/common"
	"TradeEngine/services"
	"fmt"
	"time"

	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func TradeHandler(rdb *redis.Client, client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract userId from context
		userId, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
			return
		}

		// Parse request body
		var tradeData common.Itrade
		if err := c.ShouldBindJSON(&tradeData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Set the userId from the token
		tradeData.UserId = userId.(string)
		tradeData.CreatedAt = time.Now()
		fmt.Println("Trade Data ...........", tradeData)
		// Call the Trade function
		err := services.Trade(context.Background(), tradeData, rdb, client)
		if err != nil {

			fmt.Println("Error ??????????????????????????????????????????????????????????", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Trade processed successfully"})
	}
}
