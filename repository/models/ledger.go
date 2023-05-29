package models

import (
	"nikwallet/repository/money"
	"time"
)

type TransactionType string

const (
	TransactionTypeAdd      TransactionType = "add"
	TransactionTypeWithdraw TransactionType = "withdraw"
	TransactionTypeTransfer TransactionType = "transfer"
)

type Ledger struct {
	ID              int          `gorm:"column:id"`
	FromUserID      int          `gorm:"column:from_user_id"`
	ToUserID        int          `gorm:"column:to_user_id"`
	Amount          *money.Money `gorm:"column:amount"`
	TransactionType string       `gorm:"column:transaction_type"`
	CreatedAt       time.Time    `gorm:"column:created_at"`
}
