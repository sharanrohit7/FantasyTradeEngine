package services

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func GetTradesByUserId(ctx context.Context, rdb *redis.Client, userId string) any {
	// Get all trades for the user from Redis
	tradesData := rdb.HGet(ctx, "UserTrades", userId)

	fmt.Println(tradesData)
	return "fe"
}
