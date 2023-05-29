package repository

import (
	"fmt"
	"nikwallet/repository/models"
)

func (db *PostgreSQL) CreateLedgerEntry(newEntry *models.Ledger) error {
	err := db.DB.Create(newEntry).Error
	if err != nil {
		return fmt.Errorf("failed to create ledger entry: %w", err)
	}
	return nil
}
