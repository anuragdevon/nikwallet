package repository

import (
	"fmt"

	"nikwallet/repository/models"
)

func (db *PostgreSQL) CreateWallet(newWallet *models.Wallet) (*models.Wallet, error) {
	err := db.DB.Create(newWallet).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}
	return newWallet, nil
}

func (db *PostgreSQL) GetWalletByUserID(userID int) (*models.Wallet, error) {
	wallet := &models.Wallet{}
	err := db.DB.Where("user_id = ?", userID).First(wallet).Error
	if err != nil {
		return nil, fmt.Errorf("no wallets found for user with ID %d", userID)
	}

	return wallet, nil
}

func (db *PostgreSQL) UpdateWallet(changedWallet *models.Wallet) (*models.Wallet, error) {
	err := db.DB.Save(changedWallet).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update wallet")
	}

	return changedWallet, nil
}
