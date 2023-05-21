package services

import (
	"nikwallet/repository/models"
	"nikwallet/repository/money"
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func TestWalletService(t *testing.T) {
	walletService := &WalletService{
		db: db.DB,
	}
	t.Run("CreateWallet method to create a valid wallet for successful user creation", func(t *testing.T) {
		newUser := &models.User{
			EmailID:  "testwallet511@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)

		wallet, err := walletService.CreateWallet(newUserID, money.INR)
		if err != nil {
			t.Fatalf("CreateWallet() error = %v, want nil", err)
		}

		if wallet == nil {
			t.Errorf("CreateWallet() did not set new wallet")
		}
	})

	t.Run("GetWalletByUserID method to return valid Wallet for valid userID", func(t *testing.T) {
		newUser := &models.User{
			EmailID:  "testw111@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)
		newWallet, _ := walletService.CreateWallet(newUserID, money.INR)

		wallet, err := db.GetWalletByUserID(newUserID)
		if err != nil {
			t.Fatalf("Failed to get wallet: %v", err)
		}

		if wallet.ID != newWallet.ID {
			t.Errorf("Expected wallet got = %v, want = %v", wallet, newWallet)
		}
		if wallet.UserID != newUserID {
			t.Errorf("Expected user ID %d, but got %d", newUserID, wallet.UserID)
		}
	})

	t.Run("GetWalletByUserID method to return error for NonExistentUser", func(t *testing.T) {
		_, err := db.GetWalletByUserID(9999)
		if err == nil {
			t.Fatal("Expected error, but got nil")
		}
		expectedErrMsg := "no wallets found for user with ID 9999"
		if err.Error() != expectedErrMsg {
			t.Fatalf("Expected error message '%s', but got '%v'", expectedErrMsg, err)
		}
	})

	t.Run("AddMoneyToWallet method to add money to empty wallet", func(t *testing.T) {
		newUser := &models.User{
			EmailID:  "testw13@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)

		_, _ = walletService.CreateWallet(newUserID, money.USD)

		initialMoney, _ := money.NewMoney(decimal.NewFromFloat(100.0), money.USD)

		err := walletService.AddMoneyToWallet(newUserID, *initialMoney)
		if err != nil {
			t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
		}
		updatedWallet, _ := db.GetWalletByUserID(newUserID)

		if !updatedWallet.Money.Equals(*initialMoney) {
			t.Errorf("AddMoneyToWallet() got = %v, want = %v", updatedWallet.Money, initialMoney)
		}

	})

	t.Run("AddMoneyToWallet method to add money to non empty wallet", func(t *testing.T) {
		newUser := &models.User{
			EmailID:  "testw15@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)

		_, _ = walletService.CreateWallet(newUserID, money.EUR)

		initialMoney, _ := money.NewMoney(decimal.NewFromFloat(100.0), money.EUR)
		err := walletService.AddMoneyToWallet(newUserID, *initialMoney)
		if err != nil {
			t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
		}

		additionalMoney, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.EUR)
		err = walletService.AddMoneyToWallet(newUserID, *additionalMoney)
		if err != nil {
			t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
		}
		updatedWallet, _ := db.GetWalletByUserID(newUserID)

		expectedMoney, _ := money.NewMoney(decimal.NewFromFloat(150.0), money.EUR)
		if !updatedWallet.Money.Equals(*expectedMoney) {
			t.Errorf("AddMoneyToWallet() got = %v, want = %v", updatedWallet.Money, initialMoney)
		}

	})

	t.Run("WithdrawMoneyFromWallet method to successfully return withdrawn money for valid input", func(t *testing.T) {
		newUser := &models.User{
			EmailID:  "testw17s@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)

		_, _ = walletService.CreateWallet(newUserID, money.INR)

		initialMoney, _ := money.NewMoney(decimal.NewFromFloat(100.0), money.INR)

		err := walletService.AddMoneyToWallet(newUserID, *initialMoney)
		if err != nil {
			t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
		}

		withdrawMoney, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.INR)

		withdrawnMoney, err := walletService.WithdrawMoneyFromWallet(newUserID, *withdrawMoney)
		if err != nil {
			t.Fatalf("WithdrawMoneyFromWallet() error = %v, want nil", err)
		}

		if !reflect.DeepEqual(withdrawnMoney, *withdrawMoney) {
			t.Errorf("WithdrawMoneyFromWallet() got = %v, want = %v", withdrawnMoney, withdrawMoney)
		}

		updatedWallet, _ := db.GetWalletByUserID(newUserID)

		expectedMoneyRemained, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.INR)
		if !updatedWallet.Money.Equals(*expectedMoneyRemained) {
			t.Errorf("AddMoneyToWallet() got = %v, want = %v", updatedWallet.Money, initialMoney)
		}

	})

	t.Run("WithdrawMoneyFromWallet to return error for not enough money in wallet", func(t *testing.T) {
		newUser := &models.User{
			EmailID:  "testw16s@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)

		_, _ = walletService.CreateWallet(newUserID, money.USD)

		initialMoney, _ := money.NewMoney(decimal.NewFromFloat(100.0), money.USD)

		err := walletService.AddMoneyToWallet(newUserID, *initialMoney)
		if err != nil {
			t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
		}

		withdrawMoney, _ := money.NewMoney(decimal.NewFromFloat(150.0), money.USD)

		_, err = walletService.WithdrawMoneyFromWallet(newUserID, *withdrawMoney)
		if err == nil {
			t.Error("WithdrawMoneyFromWallet() error = nil, want an error")
		}
	})

	t.Run("TrasferMoney method to successfully transfer money from sender to reiever having same currency", func(t *testing.T) {
		sender := &models.User{
			EmailID:  "test_sender@example.com",
			Password: "test123",
		}
		senderID, _ := db.CreateUser(sender)
		_, _ = walletService.CreateWallet(senderID, money.EUR)

		recipient := &models.User{
			EmailID:  "test_recipient@example.com",
			Password: "test123",
		}
		recipientID, _ := db.CreateUser(recipient)
		_, _ = walletService.CreateWallet(recipientID, money.EUR)

		initialMoney, _ := money.NewMoney(decimal.NewFromFloat(100.0), money.EUR)
		_ = walletService.AddMoneyToWallet(senderID, *initialMoney)

		transferAmount, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.EUR)
		err := walletService.TransferMoney(senderID, recipient.EmailID, *transferAmount)
		if err != nil {
			t.Fatalf("TransferMoney() error = %v, want nil", err)
		}

		expectedSenderMoney, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.EUR)
		senderWallet, _ := db.GetWalletByUserID(senderID)
		if !senderWallet.Money.Equals(*expectedSenderMoney) {
			t.Errorf("AddMoneyToWallet() got = %v, want = %v", senderWallet.Money, expectedSenderMoney)
		}

		expectedRecipientMoney, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.EUR)
		recipientWallet, _ := db.GetWalletByUserID(senderID)
		if !recipientWallet.Money.Equals(*expectedRecipientMoney) {
			t.Errorf("AddMoneyToWallet() got = %v, want = %v", recipientWallet.Money, expectedRecipientMoney)
		}

	})

	t.Run("TransferMoney method to successfully transfer money from sender to receiver having different currencies", func(t *testing.T) {
		sender := &models.User{
			EmailID:  "test_sender222@example.com",
			Password: "test123",
		}
		senderID, _ := db.CreateUser(sender)
		_, _ = walletService.CreateWallet(senderID, money.EUR)

		recipient := &models.User{
			EmailID:  "test_recipient322@example.com",
			Password: "test123",
		}
		recipientID, _ := db.CreateUser(recipient)
		_, _ = walletService.CreateWallet(recipientID, money.USD)

		initialMoney, _ := money.NewMoney(decimal.NewFromFloat(100.0), money.EUR)
		_ = walletService.AddMoneyToWallet(senderID, *initialMoney)

		transferAmount, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.EUR)
		err := walletService.TransferMoney(senderID, recipient.EmailID, *transferAmount)
		if err != nil {
			t.Fatalf("TransferMoney() error: %v, want nil", err)
		}

		expectedSenderMoney, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.EUR)
		senderWallet, _ := db.GetWalletByUserID(senderID)
		if !senderWallet.Money.Equals(*expectedSenderMoney) {
			t.Errorf("AddMoneyToWallet() got = %v, want = %v", senderWallet.Money, expectedSenderMoney)
		}

		expectedRecipientMoney, _ := money.NewMoney(decimal.NewFromFloat(45.83), money.USD)
		recipientWallet, _ := db.GetWalletByUserID(recipientID)
		if !recipientWallet.Money.Equals(*expectedRecipientMoney) {
			t.Errorf("AddMoneyToWallet() got = %v, want = %v", recipientWallet.Money, expectedRecipientMoney)
		}
	})

	t.Run("TransferMoney method to successfully transfer money from sender to receiver having different currencies", func(t *testing.T) {
		sender := &models.User{
			EmailID:  "test_sender333@example.com",
			Password: "test123",
		}
		senderID, _ := db.CreateUser(sender)
		_, _ = walletService.CreateWallet(senderID, money.EUR)

		wrongRecipientEmail := "wrong_recipient@example.com"
		wrongTransferAmount, _ := money.NewMoney(decimal.NewFromFloat(20.0), money.EUR)
		err := walletService.TransferMoney(senderID, wrongRecipientEmail, *wrongTransferAmount)
		if err == nil {
			t.Fatal("TransferMoney() expected error, got nil")
		}

	})

	t.Run("TransferMoney method should return error for invalid sender ID", func(t *testing.T) {
		recipient := &models.User{
			EmailID:  "test_recipient@example.com",
			Password: "test123",
		}
		recipientID, _ := db.CreateUser(recipient)
		_, _ = walletService.CreateWallet(recipientID, money.INR)

		transferAmount, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.INR)
		err := walletService.TransferMoney(9999, recipient.EmailID, *transferAmount)

		if err == nil {
			t.Errorf("TransferMoney() expected error but got nil for invalid sender ID")
		}
	})

	t.Run("TransferMoney method should return error for invalid recipient emailID", func(t *testing.T) {
		sender := &models.User{
			EmailID:  "test_sender@example.com",
			Password: "test123",
		}
		senderID, _ := db.CreateUser(sender)
		_, _ = walletService.CreateWallet(senderID, money.INR)

		initialMoney, _ := money.NewMoney(decimal.NewFromFloat(100.0), money.INR)
		_ = walletService.AddMoneyToWallet(senderID, *initialMoney)

		transferAmount, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.INR)
		err := walletService.TransferMoney(senderID, "invalid_recipient@example.com", *transferAmount)

		if err == nil {
			t.Errorf("TransferMoney() expected error but got nil for invalid recipient emailID")
		}
	})
}
