package repository

import (
	"nikwallet/repository/models"
	"nikwallet/repository/money"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
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
		assert.NoError(t, err)
	})

	t.Run("GetLatestLedgerEntry method to retrieve the latest ledger entry", func(t *testing.T) {
		transferAmount := &money.Money{Amount: decimal.NewFromFloat(100.0), Currency: money.INR}
		newEntry := &models.Ledger{
			SenderUserID:    1,
			ReceiverUserID:  2,
			Amount:          transferAmount,
			TransactionType: string(models.TransactionTypeTransfer),
			CreatedAt:       time.Now(),
		}
		db.CreateLedgerEntry(newEntry)

		userID := 1
		ledger, err := db.GetLatestLedgerEntry(userID)

		assert.NoError(t, err)
		assert.NotNil(t, ledger)

		assert.Equal(t, userID, ledger.SenderUserID)
	})

	t.Run("GetLastNLedgerEntries method to retrieve the last N ledger entries", func(t *testing.T) {
		userID := 1
		numEntries := 4

		for i := 1; i <= 5; i++ {
			transferAmount := &money.Money{Amount: decimal.NewFromFloat(float64(i)), Currency: money.INR}
			newEntry := &models.Ledger{
				SenderUserID:    userID,
				ReceiverUserID:  2,
				Amount:          transferAmount,
				TransactionType: string(models.TransactionTypeTransfer),
				CreatedAt:       time.Now().Add(time.Duration(-i) * time.Minute), // Create entries with decreasing timestamps
			}
			err := db.CreateLedgerEntry(newEntry)
			assert.NoError(t, err)
		}

		entries, err := db.GetLastNLedgerEntries(userID, numEntries)

		assert.NoError(t, err)
		assert.NotEmpty(t, entries)

		assert.Len(t, entries, numEntries)

		for _, entry := range entries {
			assert.Equal(t, userID, entry.SenderUserID)
		}
	})
}
