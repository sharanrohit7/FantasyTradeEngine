package services

import (
	"TradeEngine/API"
	"TradeEngine/common"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

// func Trade(ctx context.Context, tradeData common.Itrade, rdb *redis.Client, client *mongo.Client) error {
// 	// Calculate the amount
// 	amount := tradeData.Price * float32(tradeData.Quantity)

// 	// Prepare debit transaction
// 	debitTransaction := common.DebitTransaction{
// 		UserId:          tradeData.UserId,
// 		Amount:          float64(amount),
// 		WalletId:        tradeData.WalletId,
// 		CreatedAt:       time.Now(),
// 		TransactionType: "debit",
// 	}

// 	// Make external API request
// 	err := handlers.MakeDebitRequest(ctx, debitTransaction) // Assume this function is in services or refactor accordingly
// 	if err != nil {
// 		return fmt.Errorf("external API request failed: %w", err)
// 	}

// 	// Store data in Redis
// 	err = storeTradeInRedis(ctx, rdb, tradeData.UserId, tradeData, amount)
// 	if err != nil {
// 		return fmt.Errorf("error storing data in Redis: %w", err)
// 	}

// 	return nil
// }

// func Trade(ctx context.Context, tradeData common.Itrade, rdb *redis.Client, client *mongo.Client) error {
// 	// Calculate the amount
// 	amount := tradeData.Price * float32(tradeData.Quantity)

// 	// Prepare debit transaction
// 	debitTransaction := common.DebitTransaction{
// 		UserId:          tradeData.UserId,
// 		Amount:          float64(amount),
// 		WalletId:        tradeData.WalletId,
// 		CreatedAt:       time.Now(),
// 		TransactionType: "debit",
// 	}

// 	fmt.Println("Debit Transaction Body ...........................", debitTransaction)
// 	// Make external API request
// 	err := API.MakeDebitRequest(ctx, debitTransaction) // Assume this function is in services or refactor accordingly
// 	if err != nil {
// 		return fmt.Errorf("external API request failed: %w", err)
// 	}

// 	// Store data in Redis
// 	err = storeTradeInRedis(ctx, rdb, tradeData.UserId, tradeData, amount)
// 	if err != nil {
// 		return fmt.Errorf("error storing data in Redis: %w", err)
// 	}

// 	return nil
// }

func Trade(ctx context.Context, tradeData common.Itrade, rdb *redis.Client, client *mongo.Client) error {
	// Calculate the amount
	amount := tradeData.Price * float32(tradeData.Quantity)

	// Prepare debit transaction
	debitTransaction := common.DebitTransaction{
		UserId:          tradeData.UserId,
		Amount:          float64(amount),
		WalletId:        tradeData.WalletId,
		CreatedAt:       time.Now(),
		TransactionType: "debit",
	}
	debitStatus := make(chan error, 1)
	go func() {
		err := API.MakeDebitRequest(ctx, debitTransaction)
		if err != nil {
			debitStatus <- fmt.Errorf("external API request failed: %w", err)
			return
		}
		debitStatus <- nil
	}()

	// Wait for debit operation to complete
	if err := <-debitStatus; err != nil {
		return err
	}

	// Store data in Redis
	err := storeTradeInRedis(ctx, rdb, tradeData.UserId, tradeData, amount)
	if err != nil {
		return fmt.Errorf("error storing data in Redis: %w", err)
	}

	return nil
}

// func storeTradeInRedis(ctx context.Context, rdb *redis.Client, userId string, tradeData common.Itrade, amount float32) error {
// 	// Serialize trade data
// 	tradeDataStr := fmt.Sprintf("QuestionId: %s, Answer: %v, Price: %f, Quantity: %d, CreatedAt: %s, Amount: %f",
// 		tradeData.QuestionId, tradeData.Answer, tradeData.Price, tradeData.Quantity, tradeData.CreatedAt.Format(time.RFC3339), amount)

// 	// Store the trade data in Redis hash
// 	err := rdb.HSet(ctx, "UserTrades", userId, tradeDataStr).Err()
// 	if err != nil {
// 		return fmt.Errorf("error storing trade data in Redis: %w", err)
// 	}
// 	return nil
// }

func storeTradeInRedis(ctx context.Context, rdb *redis.Client, userId string, tradeData common.Itrade, amount float32) error {
	// Serialize trade data
	tradeDataStr := fmt.Sprintf("QuestionId: %s, Answer: %v, Price: %f, Quantity: %d, CreatedAt: %s, Amount: %f",
		tradeData.QuestionId, tradeData.Answer, tradeData.Price, tradeData.Quantity, tradeData.CreatedAt.Format(time.RFC3339), amount)

	// Debugging
	log.Println("Storing trade data in Redis:", tradeDataStr)

	// Store the trade data in Redis hash
	err := rdb.HSet(ctx, "UserTrades", userId, tradeDataStr).Err()
	if err != nil {
		return fmt.Errorf("error storing trade data in Redis: %w", err)
	}
	return nil
}
