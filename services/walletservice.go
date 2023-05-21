package services

import (
	"fmt"
	"nikwallet/repository"
	"nikwallet/repository/models"
	"nikwallet/repository/money"
	"time"

	"gorm.io/gorm"
)

type WalletService struct {
	db *gorm.DB
}

func NewWalletService(db *gorm.DB) *WalletService {
	return &WalletService{db: db}
}

func (ws *WalletService) CreateWallet(userID int, currency money.Currency) (*models.Wallet, error) {
	db := repository.PostgreSQL{DB: ws.db}
	_, err := db.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	initialZeroMoney, _ := money.NewMoney(money.ZeroAmountValue, currency)
	newWallet := &models.Wallet{
		UserID:    userID,
		Money:     initialZeroMoney,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return db.CreateWallet(newWallet)
}

func (ws *WalletService) GetWalletByUserID(userID int) (*models.Wallet, error) {
	db := repository.PostgreSQL{DB: ws.db}
	return db.GetWalletByUserID(userID)
}

func (ws *WalletService) AddMoneyToWallet(userID int, moneyToAdd money.Money) error {
	db := repository.PostgreSQL{DB: ws.db}

	wallet, err := db.GetWalletByUserID(userID)
	if err != nil {
		return err
	}

	newMoney, err := wallet.Money.Add(&moneyToAdd)
	if err != nil {
		return err
	}

	wallet.Money = newMoney

	err = db.UpdateWallet(wallet)
	if err != nil {
		return fmt.Errorf("failed to add money")
	}
	return nil
}

func (ws *WalletService) WithdrawMoneyFromWallet(userID int, moneyToWithdraw money.Money) (money.Money, error) {
	db := repository.PostgreSQL{DB: ws.db}
	wallet, err := db.GetWalletByUserID(userID)
	if err != nil {
		return money.Money{}, err
	}

	remainedMoney, err := wallet.Money.Subtract(&moneyToWithdraw)
	if err != nil {
		return money.Money{}, err
	}

	wallet.Money = remainedMoney
	err = db.UpdateWallet(wallet)
	if err != nil {
		return money.Money{}, fmt.Errorf("failed to withdraw money")
	}
	return moneyToWithdraw, nil
}

func (ws *WalletService) TransferMoney(senderUserID int, recipientEmail string, moneyToTransfer money.Money) error {
	db := repository.PostgreSQL{DB: ws.db}
	walletService := &WalletService{
		db: db.DB,
	}
	recipient, err := db.GetUserByEmail(recipientEmail)
	if err != nil {
		return err
	}

	recipientWallet, err := db.GetWalletByUserID(int(recipient.ID))
	if err != nil {
		return err
	}

	amountDeducted, err := walletService.WithdrawMoneyFromWallet(senderUserID, moneyToTransfer)
	if err != nil {
		return err
	}

	err = walletService.AddMoneyToWallet(recipientWallet.UserID, amountDeducted)
	if err != nil {
		_ = walletService.AddMoneyToWallet(senderUserID, amountDeducted)
		return err
	}

	return nil
}
