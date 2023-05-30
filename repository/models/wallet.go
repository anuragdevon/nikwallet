package models

import (
	"nikwallet/repository/money"
	"time"
	// pq "github.com/lib/pq"
)

type Wallet struct {
	ID     int          `gorm:"column:id"`
	UserID int          `gorm:"column:user_id"`
	Money  *money.Money `gorm:"column:amount"`
	// LedgerEntryIDs pq.Int64Array `gorm:"type:integer[]"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
