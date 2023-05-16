package wallet

import (
	"nikwallet/pkg/db"
	"nikwallet/pkg/money"
	user "nikwallet/pkg/user"
	"reflect"
	"testing"
)

func TestCreateWalletToCreateAValidWallet(t *testing.T) {
	database, err := db.ConnectToDB("testdb")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	newUser := &user.User{
		EmailID:  "testw51@example.com",
		Password: "test123",
	}
	newUserID, _ := user.CreateUser(database, newUser)

	walletID, err := CreateWallet(database, newUserID)
	if err != nil {
		t.Fatalf("CreateWallet() error = %v, want nil", err)
	}

	if walletID == 0 {
		t.Errorf("CreateWallet() did not set wallet ID")
	}
}

func TestGetWalletByIDToReturnValidWallet(t *testing.T) {
	database, err := db.ConnectToDB("testdb")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	newUser := &user.User{
		EmailID:  "testw11@example.com",
		Password: "test123",
	}
	newUserID, _ := user.CreateUser(database, newUser)

	walletID, _ := CreateWallet(database, newUserID)

	_, err = GetWalletByID(database, walletID)
	if err != nil {
		t.Fatalf("Failed to get wallet: %v", err)
	}
}

func TestAddMoneyToEmptyWallet(t *testing.T) {
	database, err := db.ConnectToDB("testdb")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	newUser := &user.User{
		EmailID:  "testw13@example.com",
		Password: "test123",
	}
	newUserID, _ := user.CreateUser(database, newUser)

	walletID, _ := CreateWallet(database, newUserID)

	initialMoney, _ := money.NewMoney(100, "INR")

	err = AddMoneyToWallet(database, walletID, *initialMoney)
	if err != nil {
		t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
	}

	wallet, err := GetWalletByID(database, walletID)
	if err != nil {
		t.Fatalf("GetWalletByID() error = %v, want nil", err)
	}

	if !reflect.DeepEqual(&wallet.Money, initialMoney) {
		t.Errorf("AddMoneyToWallet() got = %v, want = %v", wallet.Money, initialMoney)
	}
}

func TestAddMoneyToNonEmptyWallet(t *testing.T) {
	database, err := db.ConnectToDB("testdb")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	newUser := &user.User{
		EmailID:  "testw15@example.com",
		Password: "test123",
	}
	newUserID, _ := user.CreateUser(database, newUser)

	walletID, _ := CreateWallet(database, newUserID)

	initialMoney, _ := money.NewMoney(100, "INR")
	err = AddMoneyToWallet(database, walletID, *initialMoney)
	if err != nil {
		t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
	}

	additionalMoney, _ := money.NewMoney(50, "INR")
	err = AddMoneyToWallet(database, walletID, *additionalMoney)
	if err != nil {
		t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
	}

	wallet, err := GetWalletByID(database, walletID)
	if err != nil {
		t.Fatalf("GetWalletByID() error = %v, want nil", err)
	}

	expectedMoney, _ := money.NewMoney(150, "INR")
	if !reflect.DeepEqual(&wallet.Money, expectedMoney) {
		t.Errorf("AddMoneyToWallet() got = %v, want = %v", wallet.Money, expectedMoney)
	}
}
