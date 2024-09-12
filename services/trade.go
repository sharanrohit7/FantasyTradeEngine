package services

import (
	"context"
	"fmt"
	"time"

	"log"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Itrade struct {
	UserId     string  `json:"user_id" validate:"required"`
	QuestionId string  `json:"questionId" validate:"required"`
	Answer     bool    `json:"answer" validate:"required"`
	Price      float32 `json:"Price" validate:"required"`
	Quantity   int     `json:"quantity" validate:"required"`
	CreatedAt  time.Time
}

func Trade(ctx context.Context, tradeData Itrade, rdb *redis.Client, client *mongo.Client) {
	// Calculate the amount
	amount := tradeData.Price * float32(tradeData.Quantity)

	// Store data in MongoDB
	collection := client.Database("mydatabase").Collection("trades")
	_, err := collection.InsertOne(ctx, tradeData)
	if err != nil {
		log.Println("Error inserting data to MongoDB:", err)
		return
	}

	// Store data in Redis
	err = storeTradeInRedis(ctx, rdb, tradeData.UserId, tradeData, amount)
	if err != nil {
		log.Println("Error storing data in Redis:", err)
	}
}

func storeTradeInRedis(ctx context.Context, rdb *redis.Client, userId string, tradeData Itrade, amount float32) error {
	// Serialize trade data
	tradeDataStr := fmt.Sprintf("QuestionId: %s, Answer: %v, Price: %f, Quantity: %d, CreatedAt: %s, Amount: %f",
		tradeData.QuestionId, tradeData.Answer, tradeData.Price, tradeData.Quantity, tradeData.CreatedAt.Format(time.RFC3339), amount)

	// Store the trade data in Redis hash
	err := rdb.HSet(ctx, "UserTrades", userId, tradeDataStr).Err()
	if err != nil {
		return err
	}

	// Store the amount separately if needed
	err = rdb.HSet(ctx, "UserTradeAmounts", userId, amount).Err()
	if err != nil {
		return err
	}

	return nil
}
