package repository

import (
	"fmt"
	"time"

	"nikwallet/repository/models"
	"nikwallet/repository/money"
)

func (db *PostgreSQL) CreateWallet(userID int, currency money.Currency) (*models.Wallet, error) {
	_, err := db.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	initialZeroMoney, _ := money.NewMoney(money.ZeroAmountValue, currency)
	wallet := &models.Wallet{
		UserID:    userID,
		Money:     initialZeroMoney,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = db.DB.Create(wallet).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}
	return wallet, nil
}

func (db *PostgreSQL) GetWalletByUserID(userID int) (*models.Wallet, error) {
	wallet := &models.Wallet{}
	err := db.DB.Where("user_id = ?", userID).First(wallet).Error
	if err != nil {
		return nil, fmt.Errorf("no wallets found for user with ID %d", userID)
	}

	return wallet, nil
}

func (db *PostgreSQL) AddMoneyToWallet(userID int, moneyToAdd money.Money) error {
	wallet, err := db.GetWalletByUserID(userID)
	if err != nil {
		return err
	}

	newMoney, err := wallet.Money.Add(&moneyToAdd)
	if err != nil {
		return err
	}

	wallet.Money = newMoney
	err = db.DB.Save(wallet).Error
	if err != nil {
		return fmt.Errorf("failed to add money to wallet: %w", err)
	}

	return nil
}

func (db *PostgreSQL) WithdrawMoneyFromWallet(userID int, moneyToWithdraw money.Money) (money.Money, error) {
	wallet, err := db.GetWalletByUserID(userID)
	if err != nil {
		return money.Money{}, err
	}

	remainedMoney, err := wallet.Money.Subtract(&moneyToWithdraw)
	if err != nil {
		return money.Money{}, err
	}

	wallet.Money = remainedMoney
	err = db.DB.Save(wallet).Error
	if err != nil {
		return money.Money{}, fmt.Errorf("failed to withdraw money from wallet: %w", err)
	}

	return moneyToWithdraw, nil
}

func (db *PostgreSQL) TransferMoney(senderWalletID int, recipientEmail string, moneyToTransfer money.Money) error {
	recipient, err := db.GetUserByEmail(recipientEmail)
	if err != nil {
		return err
	}

	recipientWallet, err := db.GetWalletByUserID(int(recipient.ID))
	if err != nil {
		return err
	}

	amountDeducted, err := db.WithdrawMoneyFromWallet(senderWalletID, moneyToTransfer)
	if err != nil {
		return err
	}

	err = db.AddMoneyToWallet(recipientWallet.UserID, amountDeducted)
	if err != nil {
		_ = db.AddMoneyToWallet(senderWalletID, amountDeducted)
		return err
	}

	return nil
}
