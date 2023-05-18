package database

import (
	"nikwallet/database/money"
	"reflect"
	"testing"
)

func TestWallet(t *testing.T) {
	db := &PostgreSQL{}
	err := db.Connect("testdb")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	t.Run("CreateWallet method to create a valid wallet for successful user creation", func(t *testing.T) {
		newUser := &User{
			EmailID:  "testw511@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)

		wallet, err := db.CreateWallet(newUserID)
		if err != nil {
			t.Fatalf("CreateWallet() error = %v, want nil", err)
		}

		if wallet == nil {
			t.Errorf("CreateWallet() did not set new wallet")
		}
	})

	t.Run("GetWalletByUserID method to return valid Wallet for valid userID", func(t *testing.T) {
		newUser := &User{
			EmailID:  "testw111@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)
		newWallet, _ := db.CreateWallet(newUserID)

		wallet, err := db.GetWalletByUserID(newUserID)
		if err != nil {
			t.Fatalf("Failed to get wallet: %v", err)
		}

		if !reflect.DeepEqual(wallet, newWallet) {
			t.Errorf("Expected wallet got = %v, want = %v", wallet, newWallet)
		}
		if wallet.UserID != newUserID {
			t.Errorf("Expected user ID %d, but got %d", newUserID, wallet.UserID)
		}
	})

	t.Run("GetWalletByUserID method to return error for NonExistentUser", func(t *testing.T) {
		_, err = db.GetWalletByUserID(9999)
		if err == nil {
			t.Fatal("Expected error, but got nil")
		}
		expectedErrMsg := "no wallets found for user with ID 9999"
		if err.Error() != expectedErrMsg {
			t.Fatalf("Expected error message '%s', but got '%v'", expectedErrMsg, err)
		}
	})

	t.Run("AddMoneyToWallet method to add money to empty wallet", func(t *testing.T) {
		newUser := &User{
			EmailID:  "testw13@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)

		_, _ = db.CreateWallet(newUserID)

		initialMoney, _ := money.NewMoney(100, "INR")

		err := db.AddMoneyToWallet(newUserID, *initialMoney)
		if err != nil {
			t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
		}
		updatedWallet, _ := db.GetWalletByUserID(newUserID)

		if !reflect.DeepEqual(&updatedWallet.Money, initialMoney) {
			t.Errorf("AddMoneyToWallet() got = %v, want = %v", updatedWallet.Money, initialMoney)
		}
	})

	t.Run("AddMoneyToWallet method to add money to non empty wallet", func(t *testing.T) {
		newUser := &User{
			EmailID:  "testw15@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)

		_, _ = db.CreateWallet(newUserID)

		initialMoney, _ := money.NewMoney(100, "INR")
		err := db.AddMoneyToWallet(newUserID, *initialMoney)
		if err != nil {
			t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
		}

		additionalMoney, _ := money.NewMoney(50, "INR")
		err = db.AddMoneyToWallet(newUserID, *additionalMoney)
		if err != nil {
			t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
		}
		updatedWallet, _ := db.GetWalletByUserID(newUserID)

		expectedMoney, _ := money.NewMoney(150, "INR")
		if !reflect.DeepEqual(&updatedWallet.Money, expectedMoney) {
			t.Errorf("AddMoneyToWallet() got = %v, want = %v", updatedWallet.Money, expectedMoney)
		}
	})

	t.Run("WithdrawMoneyFromWallet method to successfully return withdrawn money for valid input", func(t *testing.T) {
		newUser := &User{
			EmailID:  "testw17s@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)

		_, _ = db.CreateWallet(newUserID)

		initialMoney, _ := money.NewMoney(100, "INR")

		err := db.AddMoneyToWallet(newUserID, *initialMoney)
		if err != nil {
			t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
		}

		withdrawMoney, _ := money.NewMoney(50, "INR")

		withdrawnMoney, err := db.WithdrawMoneyFromWallet(newUserID, *withdrawMoney)
		if err != nil {
			t.Fatalf("WithdrawMoneyFromWallet() error = %v, want nil", err)
		}

		if !reflect.DeepEqual(withdrawnMoney, *withdrawMoney) {
			t.Errorf("WithdrawMoneyFromWallet() got = %v, want = %v", withdrawnMoney, withdrawMoney)
		}

		updatedWallet, _ := db.GetWalletByUserID(newUserID)

		expectedMoneyRemained, _ := money.NewMoney(50, "INR")
		if !reflect.DeepEqual(&updatedWallet.Money, expectedMoneyRemained) {
			t.Errorf("WithdrawMoneyFromWallet() remained got = %v, want = %v", updatedWallet.Money, expectedMoneyRemained)
		}

	})

	t.Run("WithdrawMoneyFromWallet to return error for not enough money in wallet", func(t *testing.T) {
		newUser := &User{
			EmailID:  "testw16s@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)

		_, _ = db.CreateWallet(newUserID)

		initialMoney, _ := money.NewMoney(100, "INR")

		err := db.AddMoneyToWallet(newUserID, *initialMoney)
		if err != nil {
			t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
		}

		withdrawMoney, _ := money.NewMoney(150, "INR")

		_, err = db.WithdrawMoneyFromWallet(newUserID, *withdrawMoney)
		if err == nil {
			t.Error("WithdrawMoneyFromWallet() error = nil, want an error")
		}
	})

	t.Run("TrasferMoney method to successfully transfer money from sender to reiever", func(t *testing.T) {
		sender := &User{
			EmailID:  "test_sender@example.com",
			Password: "test123",
		}
		senderID, _ := db.CreateUser(sender)
		_, _ = db.CreateWallet(senderID)

		recipient := &User{
			EmailID:  "test_recipient@example.com",
			Password: "test123",
		}
		recipientID, _ := db.CreateUser(recipient)
		_, _ = db.CreateWallet(recipientID)

		initialMoney, _ := money.NewMoney(100, "INR")
		_ = db.AddMoneyToWallet(senderID, *initialMoney)

		transferAmount, _ := money.NewMoney(50, "INR")
		err := db.TransferMoney(senderID, recipient.EmailID, *transferAmount)
		if err != nil {
			t.Fatalf("TransferMoney() error = %v, want nil", err)
		}

		expectedSenderMoney, _ := money.NewMoney(50, "INR")
		senderWallet, _ := db.GetWalletByUserID(senderID)
		if !reflect.DeepEqual(&senderWallet.Money, expectedSenderMoney) {
			t.Errorf("TransferMoney() sender balance got = %v, want = %v", senderWallet.Money, expectedSenderMoney)
		}

		expectedRecipientMoney, _ := money.NewMoney(50, "INR")
		recipientWallet, _ := db.GetWalletByUserID(senderID)
		if !reflect.DeepEqual(&recipientWallet.Money, expectedRecipientMoney) {
			t.Errorf("TransferMoney() recipient balance got = %v, want = %v", recipientWallet.Money, expectedRecipientMoney)
		}
	})

	t.Run("TransferMoney method should return error for invalid sender ID", func(t *testing.T) {
		recipient := &User{
			EmailID:  "test_recipient@example.com",
			Password: "test123",
		}
		recipientID, _ := db.CreateUser(recipient)
		_, _ = db.CreateWallet(recipientID)

		transferAmount, _ := money.NewMoney(50, "INR")
		err := db.TransferMoney(9999, recipient.EmailID, *transferAmount)

		if err == nil {
			t.Errorf("TransferMoney() expected error but got nil for invalid sender ID")
		}
	})

	t.Run("TransferMoney method should return error for invalid recipient ID", func(t *testing.T) {
		sender := &User{
			EmailID:  "test_sender@example.com",
			Password: "test123",
		}
		senderID, _ := db.CreateUser(sender)
		_, _ = db.CreateWallet(senderID)

		initialMoney, _ := money.NewMoney(100, "INR")
		_ = db.AddMoneyToWallet(senderID, *initialMoney)

		transferAmount, _ := money.NewMoney(50, "INR")
		err := db.TransferMoney(senderID, "invalid_recipient@example.com", *transferAmount)

		if err == nil {
			t.Errorf("TransferMoney() expected error but got nil for invalid recipient ID")
		}
	})
}
