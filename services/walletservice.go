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

func (ws *WalletService) AddMoneyToWallet(userID int, moneyToAdd money.Money) (*models.Wallet, error) {
	db := repository.PostgreSQL{DB: ws.db}

	wallet, err := db.GetWalletByUserID(userID)
	if err != nil {
		return nil, err
	}

	newMoney, err := wallet.Money.Add(&moneyToAdd)
	if err != nil {
		return nil, err
	}

	wallet.Money = newMoney

	updatedWallet, err := db.UpdateWallet(wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to add money")
	}

	ledgerEntry := &models.Ledger{
		SenderUserID:    userID,
		ReceiverUserID:  userID,
		Amount:          &moneyToAdd,
		TransactionType: string(models.TransactionTypeAdd),
		CreatedAt:       time.Now(),
	}

	err = db.CreateLedgerEntry(ledgerEntry)
	if err != nil {
		return updatedWallet, fmt.Errorf("failed to create ledger entry")
	}
	return updatedWallet, nil
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
	_, err = db.UpdateWallet(wallet)
	if err != nil {
		return money.Money{}, fmt.Errorf("failed to withdraw money")
	}

	ledgerEntry := &models.Ledger{
		SenderUserID:    userID,
		ReceiverUserID:  userID,
		Amount:          &moneyToWithdraw,
		TransactionType: string(models.TransactionTypeWithdraw),
		CreatedAt:       time.Now(),
	}

	err = db.CreateLedgerEntry(ledgerEntry)
	if err != nil {
		return moneyToWithdraw, fmt.Errorf("failed to create ledger entry")
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

	_, err = walletService.AddMoneyToWallet(recipientWallet.UserID, amountDeducted)
	if err != nil {
		_, _ = walletService.AddMoneyToWallet(senderUserID, amountDeducted)
		return err
	}

	ledgerEntry := &models.Ledger{
		SenderUserID:    senderUserID,
		ReceiverUserID:  int(recipient.ID),
		Amount:          &moneyToTransfer,
		TransactionType: string(models.TransactionTypeTransfer),
		CreatedAt:       time.Now(),
	}

	err = db.CreateLedgerEntry(ledgerEntry)
	if err != nil {
		return fmt.Errorf("failed to create ledger entry")
	}

	return nil
}


func (ws *WalletService) GetLastNLedgerEntries(userID, limit int) ([]*models.Ledger, error) {
	db := repository.PostgreSQL{DB: ws.db}
	return db.GetLastNLedgerEntries(userID, limit)
}
