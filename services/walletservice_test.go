package services

import (
	"nikwallet/repository/models"
	"nikwallet/repository/money"
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
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

		_, err := walletService.AddMoneyToWallet(newUserID, *initialMoney)
		assert.NoError(t, err)

		updatedWallet, _ := db.GetWalletByUserID(newUserID)
		assert.True(t, updatedWallet.Money.Equals(*initialMoney))

		ledgerEntry, err := db.GetLatestLedgerEntry(newUserID)
		assert.NoError(t, err)

		assert.Equal(t, newUserID, ledgerEntry.SenderUserID)
		assert.Equal(t, newUserID, ledgerEntry.ReceiverUserID)
		assert.True(t, ledgerEntry.Amount.Equals(*initialMoney))
		assert.Equal(t, string(models.TransactionTypeAdd), ledgerEntry.TransactionType)
	})

	t.Run("AddMoneyToWallet method to add money to non empty wallet", func(t *testing.T) {
		newUser := &models.User{
			EmailID:  "testw15@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)

		_, _ = walletService.CreateWallet(newUserID, money.EUR)

		initialMoney, _ := money.NewMoney(decimal.NewFromFloat(100.0), money.EUR)
		_, err := walletService.AddMoneyToWallet(newUserID, *initialMoney)
		assert.Nil(t, err, "AddMoneyToWallet should not return an error")

		additionalMoney, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.EUR)
		_, err = walletService.AddMoneyToWallet(newUserID, *additionalMoney)
		assert.Nil(t, err, "AddMoneyToWallet should not return an error")

		updatedWallet, _ := db.GetWalletByUserID(newUserID)

		expectedMoney, _ := money.NewMoney(decimal.NewFromFloat(150.0), money.EUR)
		assert.True(t, updatedWallet.Money.Equals(*expectedMoney), "Wallet money should be equal to expected amount")

		latestEntry, _ := db.GetLatestLedgerEntry(newUserID)
		assert.NotNil(t, latestEntry, "Latest ledger entry should not be nil")
		assert.Equal(t, newUserID, latestEntry.SenderUserID, "SenderUserID should match the user ID")
		assert.Equal(t, newUserID, latestEntry.ReceiverUserID, "ReceiverUserID should match the user ID")
		// assert.Equal(t, *additionalMoney, *latestEntry.Amount, "Amount should match the additional money")
		assert.Equal(t, string(models.TransactionTypeAdd), latestEntry.TransactionType, "TransactionType should be 'add'")
	})

	t.Run("WithdrawMoneyFromWallet method to successfully return withdrawn money for valid input", func(t *testing.T) {
		newUser := &models.User{
			EmailID:  "testw17s@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)

		_, _ = walletService.CreateWallet(newUserID, money.INR)

		initialMoney, _ := money.NewMoney(decimal.NewFromFloat(100.0), money.INR)

		_, err := walletService.AddMoneyToWallet(newUserID, *initialMoney)
		assert.Nil(t, err, "AddMoneyToWallet should not return an error")

		withdrawMoney, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.INR)

		withdrawnMoney, err := walletService.WithdrawMoneyFromWallet(newUserID, *withdrawMoney)
		assert.Nil(t, err, "WithdrawMoneyFromWallet should not return an error")
		assert.True(t, reflect.DeepEqual(withdrawnMoney, *withdrawMoney), "Withdrawn money should be equal to the requested amount")

		updatedWallet, _ := db.GetWalletByUserID(newUserID)

		expectedMoneyRemained, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.INR)
		assert.True(t, updatedWallet.Money.Equals(*expectedMoneyRemained), "Wallet money should be equal to the expected amount")

		latestEntry, _ := db.GetLatestLedgerEntry(newUserID)
		assert.NotNil(t, latestEntry, "Latest ledger entry should not be nil")
		assert.Equal(t, newUserID, latestEntry.SenderUserID, "SenderUserID should match the user ID")
		assert.Equal(t, newUserID, latestEntry.ReceiverUserID, "ReceiverUserID should match the user ID")
		// assert.Equal(t, *withdrawMoney, *latestEntry.Amount, "Amount should match the withdrawn money")
		assert.Equal(t, string(models.TransactionTypeWithdraw), latestEntry.TransactionType, "TransactionType should be 'withdraw'")
	})

	t.Run("WithdrawMoneyFromWallet to return error for not enough money in wallet", func(t *testing.T) {
		newUser := &models.User{
			EmailID:  "testw16s@example.com",
			Password: "test123",
		}
		newUserID, _ := db.CreateUser(newUser)

		_, _ = walletService.CreateWallet(newUserID, money.USD)

		initialMoney, _ := money.NewMoney(decimal.NewFromFloat(100.0), money.USD)

		_, err := walletService.AddMoneyToWallet(newUserID, *initialMoney)
		assert.Nil(t, err, "AddMoneyToWallet should not return an error")

		withdrawMoney, _ := money.NewMoney(decimal.NewFromFloat(150.0), money.USD)

		_, err = walletService.WithdrawMoneyFromWallet(newUserID, *withdrawMoney)
		assert.NotNil(t, err, "WithdrawMoneyFromWallet should return an error")

		updatedWallet, _ := db.GetWalletByUserID(newUserID)
		assert.True(t, updatedWallet.Money.Equals(*initialMoney), "Wallet money should remain unchanged")
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
		walletService.AddMoneyToWallet(senderID, *initialMoney)

		transferAmount, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.EUR)
		err := walletService.TransferMoney(senderID, recipient.EmailID, *transferAmount)
		assert.Nil(t, err, "TransferMoney should not return an error")

		expectedSenderMoney, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.EUR)
		senderWallet, _ := db.GetWalletByUserID(senderID)
		assert.True(t, senderWallet.Money.Equals(*expectedSenderMoney), "Sender wallet money should be updated")

		expectedRecipientMoney, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.EUR)
		recipientWallet, _ := db.GetWalletByUserID(recipientID)
		assert.True(t, recipientWallet.Money.Equals(*expectedRecipientMoney), "Recipient wallet money should be updated")

		senderLedger, _ := db.GetLatestLedgerEntry(senderID)
		assert.NotNil(t, senderLedger, "Sender ledger entry should exist")
		assert.Equal(t, senderID, senderLedger.SenderUserID, "SenderUserID in ledger entry should match senderID")
		assert.Equal(t, recipientID, senderLedger.ReceiverUserID, "ReceiverUserID in ledger entry should match recipientID")
		assert.True(t, senderLedger.Amount.Equals(*transferAmount), "Amount in sender ledger entry should match transferAmount")
		assert.Equal(t, string(models.TransactionTypeTransfer), senderLedger.TransactionType, "TransactionType in sender ledger entry should be 'transfer'")

		recipientLedger, _ := db.GetLatestLedgerEntry(recipientID)
		assert.NotNil(t, recipientLedger, "Recipient ledger entry should exist")
		assert.Equal(t, senderID, recipientLedger.SenderUserID, "SenderUserID in recipient ledger entry should match senderID")
		assert.Equal(t, recipientID, recipientLedger.ReceiverUserID, "ReceiverUserID in recipient ledger entry should match recipientID")
		assert.True(t, recipientLedger.Amount.Equals(*transferAmount), "Amount in recipient ledger entry should match transferAmount")
		assert.Equal(t, string(models.TransactionTypeTransfer), recipientLedger.TransactionType, "TransactionType in recipient ledger entry should be 'transfer'")
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
		walletService.AddMoneyToWallet(senderID, *initialMoney)

		transferAmount, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.EUR)
		err := walletService.TransferMoney(senderID, recipient.EmailID, *transferAmount)
		assert.Nil(t, err, "TransferMoney should not return an error")

		expectedSenderMoney, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.EUR)
		senderWallet, _ := db.GetWalletByUserID(senderID)
		assert.True(t, senderWallet.Money.Equals(*expectedSenderMoney), "Sender wallet money should be updated")

		expectedRecipientMoney, _ := money.NewMoney(decimal.NewFromFloat(45.83), money.USD)
		recipientWallet, _ := db.GetWalletByUserID(recipientID)
		assert.True(t, recipientWallet.Money.Equals(*expectedRecipientMoney), "Recipient wallet money should be updated")

		senderLedger, _ := db.GetLatestLedgerEntry(senderID)
		assert.NotNil(t, senderLedger, "Sender ledger entry should exist")
		assert.Equal(t, senderID, senderLedger.SenderUserID, "SenderUserID in ledger entry should match senderID")
		assert.Equal(t, recipientID, senderLedger.ReceiverUserID, "ReceiverUserID in ledger entry should match recipientID")
		assert.True(t, senderLedger.Amount.Equals(*transferAmount), "Amount in sender ledger entry should match transferAmount")
		assert.Equal(t, string(models.TransactionTypeTransfer), senderLedger.TransactionType, "TransactionType in sender ledger entry should be 'transfer'")

		recipientLedger, _ := db.GetLatestLedgerEntry(recipientID)
		assert.NotNil(t, recipientLedger, "Recipient ledger entry should exist")
		assert.Equal(t, senderID, recipientLedger.SenderUserID, "SenderUserID in recipient ledger entry should match senderID")
		assert.Equal(t, recipientID, recipientLedger.ReceiverUserID, "ReceiverUserID in recipient ledger entry should match recipientID")
		// assert.True(t, recipientLedger.Amount.Equals(*expectedRecipientMoney), "Amount in recipient ledger entry should match expectedRecipientMoney")
		assert.Equal(t, string(models.TransactionTypeTransfer), recipientLedger.TransactionType, "TransactionType in recipient ledger entry should be 'transfer'")
	})

	t.Run("TransferMoney method should return error for invalid receiver ID", func(t *testing.T) {
		sender := &models.User{
			EmailID:  "test_sender333@example.com",
			Password: "test123",
		}
		senderID, _ := db.CreateUser(sender)
		_, _ = walletService.CreateWallet(senderID, money.EUR)

		wrongRecipientEmail := "wrong_recipient@example.com"
		wrongTransferAmount, _ := money.NewMoney(decimal.NewFromFloat(20.0), money.EUR)
		err := walletService.TransferMoney(senderID, wrongRecipientEmail, *wrongTransferAmount)
		assert.Error(t, err, "TransferMoney should return an error")

		// senderWallet, _ := db.GetWalletByUserID(senderID)
		// assert.Nil(t, senderWallet, "Sender wallet should not be affected")

		// ledgerEntries, _ := db.GetLatestLedgerEntry(senderID)
		// assert.Empty(t, ledgerEntries, "No ledger entries should be created")
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

		assert.Error(t, err, "TransferMoney should return an error")

		recipientWallet, _ := db.GetWalletByUserID(recipientID)
		assert.Nil(t, recipientWallet, "Recipient wallet should not be affected")

		ledgerEntries, _ := db.GetLatestLedgerEntry(recipientID)
		assert.Empty(t, ledgerEntries, "No ledger entries should be created")
	})

	t.Run("TransferMoney method should return error for invalid recipient emailID", func(t *testing.T) {
		sender := &models.User{
			EmailID:  "test_sender@example.com",
			Password: "test123",
		}
		senderID, _ := db.CreateUser(sender)
		walletService.CreateWallet(senderID, money.INR)

		initialMoney, _ := money.NewMoney(decimal.NewFromFloat(100.0), money.INR)
		walletService.AddMoneyToWallet(senderID, *initialMoney)

		transferAmount, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.INR)
		err := walletService.TransferMoney(senderID, "invalid_recipient@example.com", *transferAmount)

		assert.Error(t, err, "TransferMoney should return an error for invalid recipient emailID")

		senderWallet, _ := db.GetWalletByUserID(senderID)
		assert.Equal(t, initialMoney, senderWallet.Money, "Sender wallet money should remain unchanged")

		ledgerEntries, _ := db.GetLatestLedgerEntry(senderID)
		assert.Empty(t, ledgerEntries, "No ledger entries should be created")
	})
}
