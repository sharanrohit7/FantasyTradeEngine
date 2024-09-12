package main

import (
	"TradeEngine/config"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client // Global MongoDB client
var rdb *redis.Client    // Global Redis client

// Initialize environment and connections
func init() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// MongoDB connection
	mongoURI := os.Getenv("MONGO_URL")
	client, err = config.ConnectMongoDB(mongoURI)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Redis connection
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),     // Redis server address
		Password: os.Getenv("REDIS_PASSWORD"), // Redis password, if any
		DB:       0,                           // use default DB
	})

	// Test Redis connection
	_, err = rdb.Ping(context.TODO()).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	fmt.Println("Successfully connected to MongoDB and Redis.")
}

func main() {
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
		fmt.Println("Disconnected from MongoDB.")
	}()
}
