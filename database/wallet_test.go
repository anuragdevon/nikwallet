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

		walletID, err := db.CreateWallet(newUserID)
		if err != nil {
			t.Fatalf("CreateWallet() error = %v, want nil", err)
		}

		if walletID == 0 {
			t.Errorf("CreateWallet() did not set wallet ID")
		}
	})

	t.Run("GetWalletByID method to return valid Wallet for valid walletID", func(t *testing.T) {
		newUser := &User{
			EmailID:  "testw11@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)

		walletID, _ := db.CreateWallet(newUserID)

		_, err := db.GetWalletByID(walletID)
		if err != nil {
			t.Fatalf("Failed to get wallet: %v", err)
		}
	})

	t.Run("AddMoneyToWallet method to add money to empty wallet", func(t *testing.T) {
		newUser := &User{
			EmailID:  "testw13@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)

		walletID, _ := db.CreateWallet(newUserID)

		initialMoney, _ := money.NewMoney(100, "INR")

		err := db.AddMoneyToWallet(walletID, *initialMoney)
		if err != nil {
			t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
		}

		wallet, err := db.GetWalletByID(walletID)
		if err != nil {
			t.Fatalf("GetWalletByID() error = %v, want nil", err)
		}

		if !reflect.DeepEqual(&wallet.Money, initialMoney) {
			t.Errorf("AddMoneyToWallet() got = %v, want = %v", wallet.Money, initialMoney)
		}
	})

	t.Run("AddMoneyToWallet method to add money to non empty wallet", func(t *testing.T) {
		newUser := &User{
			EmailID:  "testw15@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)

		walletID, _ := db.CreateWallet(newUserID)

		initialMoney, _ := money.NewMoney(100, "INR")
		err := db.AddMoneyToWallet(walletID, *initialMoney)
		if err != nil {
			t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
		}

		additionalMoney, _ := money.NewMoney(50, "INR")
		err = db.AddMoneyToWallet(walletID, *additionalMoney)
		if err != nil {
			t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
		}

		wallet, err := db.GetWalletByID(walletID)
		if err != nil {
			t.Fatalf("GetWalletByID() error = %v, want nil", err)
		}

		expectedMoney, _ := money.NewMoney(150, "INR")
		if !reflect.DeepEqual(&wallet.Money, expectedMoney) {
			t.Errorf("AddMoneyToWallet() got = %v, want = %v", wallet.Money, expectedMoney)
		}
	})

	t.Run("WithdrawMoneyFromWallet method to successfully return withdrawn money for valid input", func(t *testing.T) {
		newUser := &User{
			EmailID:  "testw17s@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)

		walletID, _ := db.CreateWallet(newUserID)

		initialMoney, _ := money.NewMoney(100, "INR")

		err := db.AddMoneyToWallet(walletID, *initialMoney)
		if err != nil {
			t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
		}

		withdrawMoney, _ := money.NewMoney(50, "INR")

		withdrawnMoney, err := db.WithdrawMoneyFromWallet(walletID, *withdrawMoney)
		if err != nil {
			t.Fatalf("WithdrawMoneyFromWallet() error = %v, want nil", err)
		}

		if !reflect.DeepEqual(withdrawnMoney, *withdrawMoney) {
			t.Errorf("WithdrawMoneyFromWallet() got = %v, want = %v", withdrawnMoney, withdrawMoney)
		}
	})

	t.Run("WithdrawMoneyFromWallet to return error for not enough money in wallet", func(t *testing.T) {
		newUser := &User{
			EmailID:  "testw16s@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)

		walletID, _ := db.CreateWallet(newUserID)

		initialMoney, _ := money.NewMoney(100, "INR")

		err := db.AddMoneyToWallet(walletID, *initialMoney)
		if err != nil {
			t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
		}

		withdrawMoney, _ := money.NewMoney(150, "INR")

		_, err = db.WithdrawMoneyFromWallet(walletID, *withdrawMoney)
		if err == nil {
			t.Error("WithdrawMoneyFromWallet() error = nil, want an error")
		}
	})
}
