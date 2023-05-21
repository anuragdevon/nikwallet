package repository

import (
	"nikwallet/repository/models"
	"nikwallet/repository/money"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestWallet(t *testing.T) {
	t.Run("CreateWallet method to create a valid wallet for successful user creation", func(t *testing.T) {
		newUser := &models.User{
			EmailID:  "test243@example.com",
			Password: "test123",
		}
		userID, err := db.CreateUser(newUser)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		newWallet := &models.Wallet{
			UserID: userID,
			Money:  &money.Money{Amount: money.ZeroAmountValue, Currency: money.INR},
		}
		createdWallet, err := db.CreateWallet(newWallet)

		assert.NoError(t, err)
		assert.NotNil(t, createdWallet)
		assert.Equal(t, 1, createdWallet.ID)
	})

	t.Run("GetWalletByUserID method to return valid Wallet for valid userID", func(t *testing.T) {
		newUser := &models.User{
			EmailID:  "testw133@example.com",
			Password: "test123",
		}
		userID, _ := db.CreateUser(newUser)
		newWallet := &models.Wallet{
			UserID: userID,
			Money:  &money.Money{Amount: money.ZeroAmountValue, Currency: money.INR},
		}
		createdWallet, _ := db.CreateWallet(newWallet)

		wallet, err := db.GetWalletByUserID(userID)
		assert.Nil(t, err)
		assert.Equal(t, createdWallet.ID, wallet.ID)
		assert.Equal(t, createdWallet.UserID, wallet.UserID)
	})

	t.Run("GetWalletByUserID method to return error for NonExistentUser", func(t *testing.T) {
		_, err := db.GetWalletByUserID(9999)
		assert.NotNil(t, err)

		expectedErrMsg := "no wallets found for user with ID 9999"
		assert.Equal(t, expectedErrMsg, err.Error())
	})

	t.Run("UpdateWallet method to update wallet successfully", func(t *testing.T) {
		newUser := &models.User{
			EmailID:  "testw122@example.com",
			Password: "test123",
		}
		userID, _ := db.CreateUser(newUser)
		newWallet := &models.Wallet{
			UserID: userID,
			Money:  &money.Money{Amount: money.ZeroAmountValue, Currency: money.INR},
		}
		wallet, _ := db.CreateWallet(newWallet)

		wallet.Money = &money.Money{Amount: decimal.NewFromFloat(100.0), Currency: money.INR}
		err := db.UpdateWallet(wallet)
		assert.Nil(t, err)
	})

}
