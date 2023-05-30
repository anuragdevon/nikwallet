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

func (db *PostgreSQL) GetLatestLedgerEntry(userID int) (*models.Ledger, error) {
	ledger := &models.Ledger{}
	err := db.DB.Where("sender_user_id = ? OR receiver_user_id = ?", userID, userID).
		Order("created_at DESC").
		First(ledger).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve latest ledger entry: %w", err)
	}
	return ledger, nil
}

func (db *PostgreSQL) GetLastNLedgerEntries(userID, limit int) ([]*models.Ledger, error) {
	var entries []*models.Ledger
	err := db.DB.Where("sender_user_id = ? OR receiver_user_id = ?", userID, userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve last N ledger entries: %w", err)
	}
	return entries, nil
}