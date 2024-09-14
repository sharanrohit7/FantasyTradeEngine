package common

import "time"

// Itrade represents a trade with necessary fields.
type Itrade struct {
	UserId     string  `json:"user_id" validate:"required"`
	QuestionId string  `json:"questionId" validate:"required"`
	WalletId   string  `json:"wallet_id" validate:"required"`
	Answer     bool    `json:"answer" validate:"required"`
	Price      float32 `json:"Price" validate:"required"`
	Quantity   int     `json:"quantity" validate:"required"`
	CreatedAt  time.Time
}

// DebitTransaction represents a debit transaction structure.
type DebitTransaction struct {
	UserId          string    `json:"user_id"`
	Amount          float64   `json:"amount"`
	WalletId        string    `json:"wallet_id"`
	CreatedAt       time.Time `json:"created_at"`
	TransactionType string    `json:"transaction_type"`
}
