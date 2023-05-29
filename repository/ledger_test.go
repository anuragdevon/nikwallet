package repository

import (
	"nikwallet/repository/models"
	"nikwallet/repository/money"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestLedger(t *testing.T) {
	t.Run("CreateLedgerEntry method to successfully creates new entry in ledger", func(t *testing.T) {
		transferAmount := &money.Money{Amount: decimal.NewFromFloat(100.0), Currency: money.INR}
		newEntry := &models.Ledger{
			SenderUserID:    1,
			ReceiverUserID:  2,
			Amount:          transferAmount,
			TransactionType: string(models.TransactionTypeTransfer),
			CreatedAt:       time.Now(),
		}
		err := db.CreateLedgerEntry(newEntry)
		if err != nil {
			t.Fatalf("failed to create new ledger entry: %v", err)
		}
	})
}
