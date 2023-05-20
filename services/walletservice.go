package services

import (
	"nikwallet/repository"
	"nikwallet/repository/models"
	"nikwallet/repository/money"

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
	return db.CreateWallet(userID, currency)
}

func (ws *WalletService) GetWalletByUserID(userID int) (*models.Wallet, error) {
	db := repository.PostgreSQL{DB: ws.db}
	return db.GetWalletByUserID(userID)
}

func (ws *WalletService) AddMoneyToWallet(walletID int, moneyToAdd money.Money) error {
	db := repository.PostgreSQL{DB: ws.db}
	return db.AddMoneyToWallet(walletID, moneyToAdd)
}

func (ws *WalletService) WithdrawMoneyFromWallet(walletID int, moneyToWithdraw money.Money) (money.Money, error) {
	db := repository.PostgreSQL{DB: ws.db}
	return db.WithdrawMoneyFromWallet(walletID, moneyToWithdraw)
}

func (ws *WalletService) TransferMoney(walletID int, recipientEmail string, moneyToTransfer money.Money) error {
	db := repository.PostgreSQL{DB: ws.db}
	return db.TransferMoney(walletID, recipientEmail, moneyToTransfer)
}
