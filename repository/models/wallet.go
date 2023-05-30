package models

import (
	"nikwallet/repository/money"
	"time"
)

type Wallet struct {
	ID     int          `gorm:"column:id"`
	UserID int          `gorm:"column:user_id"`
	Money  *money.Money `gorm:"column:amount"`
	// LedgerEntryIDs []int64      `gorm:"column:ledger_entry_ids"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
