package services

import (
	"nikwallet/pkg/money"
	"nikwallet/pkg/wallet"

	"database/sql"
)

type WalletService struct {
	db *sql.DB
}

func NewWalletService(db *sql.DB) *WalletService {
	return &WalletService{db}
}

func (s *WalletService) CreateWallet(userID int) (int, error) {
	return wallet.CreateWallet(userID)
}

func (s *WalletService) GetWalletByID(walletID int) (*wallet.Wallet, error) {
	return wallet.GetWalletByID(walletID)
}

func (s *WalletService) AddMoneyToWallet(walletID int, moneyToAdd money.Money) error {
	return wallet.AddMoneyToWallet(walletID, moneyToAdd)
}

func (s *WalletService) WithdrawMoneyFromWallet(walletID int, moneyToWithdraw money.Money) (money.Money, error) {
	return wallet.WithdrawMoneyFromWallet(walletID, moneyToWithdraw)
}
