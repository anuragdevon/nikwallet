package wallet

import (
	"nikwallet/pkg/db"
	user "nikwallet/pkg/user"
	"testing"
)

func TestCreateWalletToCreateAValidWallet(t *testing.T) {
	database, err := db.ConnectToDB("testdb")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	newUser := &user.User{
		EmailID:  "testw5@example.com",
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
