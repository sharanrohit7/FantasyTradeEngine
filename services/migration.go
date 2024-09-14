package services

import (
	"TradeEngine/common"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MigrateRedisToMongo migrates trade data from Redis to MongoDB, avoiding duplicates
func MigrateRedisToMongo(ctx context.Context, rdb *redis.Client, client *mongo.Client) error {
	// Define the Redis key for trades
	tradesKey := "UserTrades"

	// Fetch all trades from Redis
	tradesMap, err := rdb.HGetAll(ctx, tradesKey).Result()
	if err != nil {
		return fmt.Errorf("error fetching trades from Redis: %w", err)
	}

	// Prepare bulk write operations for MongoDB
	var bulkWrites []mongo.WriteModel
	for userId, tradeDataStr := range tradesMap {
		// Parse trade data from Redis string
		tradeData, err := parseTradeData(tradeDataStr)
		if err != nil {
			return fmt.Errorf("error parsing trade data for user %s: %w", userId, err)
		}

		// Set the UserId field (since it's not part of the serialized string)
		tradeData.UserId = userId

		// Prepare the MongoDB filter and document
		filter := bson.M{
			"user_id":     tradeData.UserId,
			"question_id": tradeData.QuestionId,
			"created_at":  tradeData.CreatedAt,
		}
		update := bson.M{
			"$set": bson.M{
				"user_id":     tradeData.UserId,
				"question_id": tradeData.QuestionId,
				"answer":      tradeData.Answer,
				"price":       tradeData.Price,
				"quantity":    tradeData.Quantity,
				"created_at":  tradeData.CreatedAt,
			},
		}

		// Use upsert to insert if the document doesn't exist or update it if it does
		bulkWrites = append(bulkWrites, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true))
	}

	// Perform bulk upsert into MongoDB
	collection := client.Database("mydatabase").Collection("trades")
	if len(bulkWrites) > 0 {
		_, err = collection.BulkWrite(ctx, bulkWrites)
		if err != nil {
			return fmt.Errorf("error bulk upserting trades into MongoDB: %w", err)
		}
	}

	log.Println("Migration from Redis to MongoDB completed successfully")
	return nil
}

// parseTradeData parses the serialized trade data string into the Itrade struct
func parseTradeData(dataStr string) (common.Itrade, error) {
	var tradeData common.Itrade
	dataStr = strings.Trim(dataStr, " ")

	// Split the trade data string by ', ' to get key-value pairs
	parts := strings.Split(dataStr, ", ")
	for _, part := range parts {
		keyValue := strings.Split(part, ": ")
		if len(keyValue) != 2 {
			return tradeData, fmt.Errorf("invalid trade data format")
		}

		key := keyValue[0]
		value := keyValue[1]

		// Match each key with its corresponding Itrade field
		switch key {
		case "QuestionId":
			tradeData.QuestionId = value
		case "Answer":
			tradeData.Answer = (value == "true")
		case "Price":
			price, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return tradeData, fmt.Errorf("error parsing Price: %w", err)
			}
			tradeData.Price = float32(price)
		case "Quantity":
			quantity, err := strconv.Atoi(value)
			if err != nil {
				return tradeData, fmt.Errorf("error parsing Quantity: %w", err)
			}
			tradeData.Quantity = quantity
		case "CreatedAt":
			createdAt, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return tradeData, fmt.Errorf("error parsing CreatedAt: %w", err)
			}
			tradeData.CreatedAt = createdAt
		}
	}

	return tradeData, nil
}
