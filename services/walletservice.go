package services

import (
	"nikwallet/database"
	"nikwallet/database/money"

	"database/sql"
)

type WalletService struct {
	db *sql.DB
}

func NewWalletService(db *sql.DB) *WalletService {
	return &WalletService{db}
}

func (ws *WalletService) CreateWallet(userID int) (int, error) {
	db := database.PostgreSQL{DB: ws.db}
	return db.CreateWallet(userID)
}

func (ws *WalletService) GetWalletByID(walletID int) (*database.Wallet, error) {
	db := database.PostgreSQL{DB: ws.db}
	return db.GetWalletByID(walletID)
}

func (ws *WalletService) AddMoneyToWallet(walletID int, moneyToAdd money.Money) error {
	db := database.PostgreSQL{DB: ws.db}
	return db.AddMoneyToWallet(walletID, moneyToAdd)
}

func (ws *WalletService) WithdrawMoneyFromWallet(walletID int, moneyToWithdraw money.Money) (money.Money, error) {
	db := database.PostgreSQL{DB: ws.db}
	return db.WithdrawMoneyFromWallet(walletID, moneyToWithdraw)
}
