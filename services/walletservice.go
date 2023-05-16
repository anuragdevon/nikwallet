package services

import (
	"nikwallet/pkg/db"
	"nikwallet/pkg/money"
	"nikwallet/pkg/user"
	"nikwallet/pkg/wallet"
)

type WalletService struct {
	db *db.DB
}

func NewWalletService(db *db.DB) *WalletService {
	return &WalletService{db: db}
}

func (s *WalletService) CreateWallet(userID int) (int, error) {
	return wallet.CreateWallet(s.db, userID)
}

func (s *WalletService) GetWalletByID(walletID int) (*wallet.Wallet, error) {
	return wallet.GetWalletByID(s.db, walletID)
}

func (s *WalletService) AddMoneyToWallet(walletID int, moneyToAdd money.Money) error {
	return wallet.AddMoneyToWallet(s.db, walletID, moneyToAdd)
}

func (s *WalletService) WithdrawMoneyFromWallet(walletID int, moneyToWithdraw money.Money) (money.Money, error) {
	return wallet.WithdrawMoneyFromWallet(s.db, walletID, moneyToWithdraw)
}

func (s *WalletService) GetUserByID(userID int) (*user.User, error) {
	return user.GetUserByID(s.db, userID)
}
