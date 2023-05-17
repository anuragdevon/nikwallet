package wallet

import (
	"database/sql"
	"log"
	"nikwallet/pkg/db"
	"nikwallet/pkg/money"
	"nikwallet/pkg/user"
	"os"
	"reflect"
	"testing"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	if err := db.ConnectToDB("testdb"); err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}
	testDB = db.DB

	code := m.Run()

	testDB.Close()

	os.Exit(code)
}

func TestCreateWalletToCreateAValidWallet(t *testing.T) {
	newUser := &user.User{
		EmailID:  "testw511@example.com",
		Password: "test123",
	}
	newUserID, _ := user.CreateUser(newUser)

	walletID, err := CreateWallet(newUserID)
	if err != nil {
		t.Fatalf("CreateWallet() error = %v, want nil", err)
	}

	if walletID == 0 {
		t.Errorf("CreateWallet() did not set wallet ID")
	}
}

func TestGetWalletByIDToReturnValidWallet(t *testing.T) {
	newUser := &user.User{
		EmailID:  "testw11@example.com",
		Password: "test123",
	}
	newUserID, _ := user.CreateUser(newUser)

	walletID, _ := CreateWallet(newUserID)

	_, err := GetWalletByID(walletID)
	if err != nil {
		t.Fatalf("Failed to get wallet: %v", err)
	}
}

func TestAddMoneyToEmptyWallet(t *testing.T) {
	newUser := &user.User{
		EmailID:  "testw13@example.com",
		Password: "test123",
	}
	newUserID, _ := user.CreateUser(newUser)

	walletID, _ := CreateWallet(newUserID)

	initialMoney, _ := money.NewMoney(100, "INR")

	err := AddMoneyToWallet(walletID, *initialMoney)
	if err != nil {
		t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
	}

	wallet, err := GetWalletByID(walletID)
	if err != nil {
		t.Fatalf("GetWalletByID() error = %v, want nil", err)
	}

	if !reflect.DeepEqual(&wallet.Money, initialMoney) {
		t.Errorf("AddMoneyToWallet() got = %v, want = %v", wallet.Money, initialMoney)
	}
}

func TestAddMoneyToNonEmptyWallet(t *testing.T) {
	newUser := &user.User{
		EmailID:  "testw15@example.com",
		Password: "test123",
	}
	newUserID, _ := user.CreateUser(newUser)

	walletID, _ := CreateWallet(newUserID)

	initialMoney, _ := money.NewMoney(100, "INR")
	err := AddMoneyToWallet(walletID, *initialMoney)
	if err != nil {
		t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
	}

	additionalMoney, _ := money.NewMoney(50, "INR")
	err = AddMoneyToWallet(walletID, *additionalMoney)
	if err != nil {
		t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
	}

	wallet, err := GetWalletByID(walletID)
	if err != nil {
		t.Fatalf("GetWalletByID() error = %v, want nil", err)
	}

	expectedMoney, _ := money.NewMoney(150, "INR")
	if !reflect.DeepEqual(&wallet.Money, expectedMoney) {
		t.Errorf("AddMoneyToWallet() got = %v, want = %v", wallet.Money, expectedMoney)
	}
}

func TestWithdrawMoneyFromWalletToReturnWithdrawnMoney(t *testing.T) {
	newUser := &user.User{
		EmailID:  "testw17s@example.com",
		Password: "test123",
	}
	newUserID, _ := user.CreateUser(newUser)

	walletID, _ := CreateWallet(newUserID)

	initialMoney, _ := money.NewMoney(100, "INR")

	err := AddMoneyToWallet(walletID, *initialMoney)
	if err != nil {
		t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
	}

	withdrawMoney, _ := money.NewMoney(50, "INR")

	withdrawnMoney, err := WithdrawMoneyFromWallet(walletID, *withdrawMoney)
	if err != nil {
		t.Fatalf("WithdrawMoneyFromWallet() error = %v, want nil", err)
	}

	if !reflect.DeepEqual(withdrawnMoney, *withdrawMoney) {
		t.Errorf("WithdrawMoneyFromWallet() got = %v, want = %v", withdrawnMoney, withdrawMoney)
	}
}

func TestWithdrawMoneyFromWalletToReturnErrorForNotEnoughMoney(t *testing.T) {
	newUser := &user.User{
		EmailID:  "testw16s@example.com",
		Password: "test123",
	}
	newUserID, _ := user.CreateUser(newUser)

	walletID, _ := CreateWallet(newUserID)

	initialMoney, _ := money.NewMoney(100, "INR")

	err := AddMoneyToWallet(walletID, *initialMoney)
	if err != nil {
		t.Fatalf("AddMoneyToWallet() error = %v, want nil", err)
	}

	withdrawMoney, _ := money.NewMoney(150, "INR")

	_, err = WithdrawMoneyFromWallet(walletID, *withdrawMoney)
	if err == nil {
		t.Error("WithdrawMoneyFromWallet() error = nil, want an error")
	}
}
