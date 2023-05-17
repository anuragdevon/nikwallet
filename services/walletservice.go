package services

import (
	"nikwallet/pkg/money"
	"nikwallet/pkg/user"
	"nikwallet/pkg/wallet"
)

type WalletService struct {
}

func NewWalletService() *WalletService {
	return &WalletService{}
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

func (s *WalletService) GetUserByID(userID int) (*user.User, error) {
	return user.GetUserByID(userID)
}
