package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"nikwallet/handlers/dto"
	"nikwallet/repository/models"
	"nikwallet/repository/money"
	"nikwallet/services"
)

func TestWalletHandlers(t *testing.T) {

	userService := services.NewUserService(db.DB)
	authService := services.NewAuthService(db.DB)
	walletService := services.NewWalletService(db.DB)

	walletHandlers := NewWalletHandlers(walletService, authService, userService)

	t.Run("CreateWalletHandler to return 201 StatusCreated for successfully create wallet", func(t *testing.T) {
		newUser := &models.User{
			EmailID:  "testw5111@example.com",
			Password: "password",
		}

		_, err := userService.CreateUser(newUser)
		assert.NoError(t, err)

		token, err := authService.AuthenticateUser(newUser.EmailID, newUser.Password)
		assert.Nil(t, err)
		assert.NotNil(t, token)

		userCurrency := map[string]interface{}{
			"currency": "INR",
		}

		reqBody, err := json.Marshal(userCurrency)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "/wallet", bytes.NewReader(reqBody))
		req.Header.Set("id_token", token)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(walletHandlers.CreateWalletHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)

		var wallet *models.Wallet
		err = json.NewDecoder(recorder.Body).Decode(&wallet)
		assert.NoError(t, err)
	})

	t.Run("AddMoneyToWalletHandler to return 200 StatusOk for successfull add money to user's wallet", func(t *testing.T) {
		newUser := &models.User{
			EmailID:  "testw5112@example.com",
			Password: "password",
		}

		userID, err := userService.CreateUser(newUser)
		assert.NoError(t, err)

		_, err = walletService.CreateWallet(userID, money.INR)
		assert.NoError(t, err)

		IDToken, _ := authService.AuthenticateUser(newUser.EmailID, newUser.Password)
		assert.NotNil(t, IDToken)

		addMoneyRequest := money.Money{Amount: decimal.NewFromFloat(50.0), Currency: money.INR}
		reqBody, err := json.Marshal(addMoneyRequest)
		assert.NoError(t, err)

		url := "/wallet"
		req, err := http.NewRequest("PUT", url, bytes.NewReader(reqBody))
		req.Header.Set("id_token", IDToken)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(walletHandlers.AddMoneyToWalletHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response dto.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, "money added to wallet successfully", response.Message)

		expectedAddedMoney, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.INR)
		senderWallet, _ := walletService.GetWalletByUserID(userID)
		if !senderWallet.Money.Equals(*expectedAddedMoney) {
			t.Errorf("AddMoneyToWallet() got = %v, want = %v", senderWallet.Money, expectedAddedMoney)
		}

	})

	t.Run("AddMoneyToWalletHandler to return status 400 bad request for InvalidAmount", func(t *testing.T) {
		newUser := &models.User{
			EmailID:  "testw599@example.com",
			Password: "password",
		}

		userID, err := userService.CreateUser(newUser)
		assert.NoError(t, err)

		_, err = walletService.CreateWallet(userID, money.INR)
		assert.NoError(t, err)

		IDToken, _ := authService.AuthenticateUser(newUser.EmailID, newUser.Password)
		assert.NotNil(t, IDToken)

		invalidWithdrawMoneyRequest := map[string]interface{}{
			"amount":   "notanumber",
			"currency": "INR",
		}

		reqBody, err := json.Marshal(invalidWithdrawMoneyRequest)
		assert.NoError(t, err)

		url := "/wallet"
		req, err := http.NewRequest("PUT", url, bytes.NewReader(reqBody))
		req.Header.Set("id_token", IDToken)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(walletHandlers.AddMoneyToWalletHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var response dto.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		assert.Error(t, err, "invalid amount")
	})

	t.Run("WithdrawMoneyFromWalletHandler to return 200 StatusOk for successfull withdraw money from user's wallet", func(t *testing.T) {
		newUser := &models.User{
			EmailID:  "testw5113@example.com",
			Password: "password",
		}

		userID, err := userService.CreateUser(newUser)
		assert.NoError(t, err)

		_, err = walletService.CreateWallet(userID, money.INR)
		assert.NoError(t, err)

		IDToken, _ := authService.AuthenticateUser(newUser.EmailID, newUser.Password)
		assert.NotNil(t, IDToken)

		addMoneyRequest := money.Money{Amount: decimal.NewFromFloat(50.0), Currency: "INR"}

		err = walletService.AddMoneyToWallet(userID, addMoneyRequest)
		assert.NoError(t, err)

		withdrawMoneyRequest := money.Money{Amount: decimal.NewFromFloat(50.0), Currency: "INR"}
		reqBody, err := json.Marshal(withdrawMoneyRequest)
		assert.NoError(t, err)

		url := "wallet/withdraw"
		req, err := http.NewRequest("PUT", url, bytes.NewReader(reqBody))
		req.Header.Set("id_token", IDToken)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(walletHandlers.WithdrawMoneyFromWalletHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var wallet *models.Wallet
		err = json.NewDecoder(recorder.Body).Decode(&wallet)
		assert.NoError(t, err)

		expectedRemainedMoney, _ := money.NewMoney(decimal.NewFromFloat(0.0), money.INR)
		senderWallet, _ := walletService.GetWalletByUserID(userID)
		if !senderWallet.Money.Equals(*expectedRemainedMoney) {
			t.Errorf("AddMoneyToWallet() got = %v, want = %v", senderWallet.Money, expectedRemainedMoney)
		}
	})

	t.Run("WithdrawMoneyFromWalletHandler to return status 400 bad request for InsufficientFunds", func(t *testing.T) {
		newUser := &models.User{
			EmailID:  "testw5114@example.com",
			Password: "password",
		}

		userID, err := userService.CreateUser(newUser)
		assert.NoError(t, err)

		_, err = walletService.CreateWallet(userID, money.INR)
		assert.NoError(t, err)

		IDToken, _ := authService.AuthenticateUser(newUser.EmailID, newUser.Password)
		assert.NotNil(t, IDToken)

		addMoneyRequest := money.Money{Amount: decimal.NewFromFloat(40.0), Currency: money.INR}
		err = walletService.AddMoneyToWallet(userID, addMoneyRequest)
		assert.NoError(t, err)

		withdrawMoneyRequest := money.Money{Amount: decimal.NewFromFloat(50.0), Currency: money.INR}
		reqBody, err := json.Marshal(withdrawMoneyRequest)
		assert.NoError(t, err)

		url := "/wallet/withdraw"
		req, err := http.NewRequest("PUT", url, bytes.NewReader(reqBody))
		req.Header.Set("id_token", IDToken)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(walletHandlers.WithdrawMoneyFromWalletHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var response dto.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		assert.Error(t, err, "insufficient funds")
	})

	t.Run("WithdrawMoneyFromWalletHandler to return status 400 BadRequest for InvalidAmount", func(t *testing.T) {
		newUser := &models.User{
			EmailID:  "testw5115@example.com",
			Password: "password",
		}

		userID, err := userService.CreateUser(newUser)
		assert.NoError(t, err)

		_, err = walletService.CreateWallet(userID, money.INR)
		assert.NoError(t, err)

		IDToken, _ := authService.AuthenticateUser(newUser.EmailID, newUser.Password)
		assert.NotNil(t, IDToken)

		invalidWithdrawMoneyRequest := map[string]interface{}{
			"amount":   "notanumber",
			"currency": "INR",
		}

		reqBody, err := json.Marshal(invalidWithdrawMoneyRequest)
		assert.NoError(t, err)

		url := "/wallet/withdraw"
		req, err := http.NewRequest("PUT", url, bytes.NewReader(reqBody))
		req.Header.Set("id_token", IDToken)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(walletHandlers.WithdrawMoneyFromWalletHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var response dto.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		assert.Error(t, err, "invalid amount")
	})

	t.Run("TransferMoneyHandler to return 200 StatusOk for successful transfer of money from sender to reciever having same currency", func(t *testing.T) {
		sender := &models.User{
			EmailID:  "sender@example.com",
			Password: "test123",
		}
		senderID, _ := userService.CreateUser(sender)
		_, _ = walletService.CreateWallet(senderID, money.INR)

		recipient := &models.User{
			EmailID:  "recipient@example.com",
			Password: "test123",
		}
		recipientID, _ := userService.CreateUser(recipient)
		_, _ = walletService.CreateWallet(recipientID, money.INR)

		initialMoney, _ := money.NewMoney(decimal.NewFromFloat(100.0), money.INR)
		_ = walletService.AddMoneyToWallet(senderID, *initialMoney)

		IDToken, _ := authService.AuthenticateUser(sender.EmailID, sender.Password)

		transferMoney, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.INR)

		transferMoneyPayload := dto.MoneyTransferDTO{
			Amount:         transferMoney,
			RecipientEmail: recipient.EmailID,
		}

		reqBody, _ := json.Marshal(transferMoneyPayload)

		url := "/wallet/transfer"
		req, _ := http.NewRequest("PUT", url, bytes.NewReader(reqBody))
		req.Header.Set("id_token", IDToken)

		recorder := httptest.NewRecorder()
		http.HandlerFunc(walletHandlers.TransferMoneyHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		var response dto.Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "money transferred successfully", response.Message)

		expectedSenderMoney, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.INR)
		senderWallet, _ := walletService.GetWalletByUserID(senderID)
		if !senderWallet.Money.Equals(*expectedSenderMoney) {
			t.Errorf("AddMoneyToWallet() got = %v, want = %v", senderWallet.Money, initialMoney)
		}

		expectedRecipientMoney, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.INR)
		recipientWallet, _ := walletService.GetWalletByUserID(recipientID)

		if !recipientWallet.Money.Equals(*expectedRecipientMoney) {
			t.Errorf("AddMoneyToWallet() got = %v, want = %v", recipientWallet.Money, expectedRecipientMoney)
		}
	})

	t.Run("TransferMoneyHandler to return 200 StatusOk for successful transfer of money from sender to reciever having different currency", func(t *testing.T) {
		sender := &models.User{
			EmailID:  "senderUSA@example.com",
			Password: "test123",
		}
		senderID, _ := userService.CreateUser(sender)
		_, _ = walletService.CreateWallet(senderID, money.USD)

		recipient := &models.User{
			EmailID:  "recipientFrance@example.com",
			Password: "test123",
		}
		recipientID, _ := userService.CreateUser(recipient)
		_, _ = walletService.CreateWallet(recipientID, money.EUR)

		initialMoney, _ := money.NewMoney(decimal.NewFromFloat(10.0), money.USD)
		_ = walletService.AddMoneyToWallet(senderID, *initialMoney)

		IDToken, _ := authService.AuthenticateUser(sender.EmailID, sender.Password)

		transferMoney, _ := money.NewMoney(decimal.NewFromFloat(2.0), money.USD)

		transferMoneyPayload := dto.MoneyTransferDTO{
			Amount:         transferMoney,
			RecipientEmail: recipient.EmailID,
		}

		reqBody, _ := json.Marshal(transferMoneyPayload)

		url := "/wallet/transfer"
		req, _ := http.NewRequest("PUT", url, bytes.NewReader(reqBody))
		req.Header.Set("id_token", IDToken)

		recorder := httptest.NewRecorder()
		http.HandlerFunc(walletHandlers.TransferMoneyHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		var response dto.Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "money transferred successfully", response.Message)

		expectedSenderMoney, _ := money.NewMoney(decimal.NewFromFloat(8.0), money.USD)
		senderWallet, _ := walletService.GetWalletByUserID(senderID)
		if !senderWallet.Money.Equals(*expectedSenderMoney) {
			t.Errorf("AddMoneyToWallet() got = %v, want = %v", senderWallet.Money, expectedSenderMoney)
		}

		expectedRecipientMoney, _ := money.NewMoney(decimal.NewFromFloat(2.18), money.EUR)
		recipientWallet, _ := walletService.GetWalletByUserID(recipientID)
		if !recipientWallet.Money.Equals(*expectedRecipientMoney) {
			t.Errorf("AddMoneyToWallet() got = %v, want = %v", recipientWallet.Money, expectedRecipientMoney)
		}
	})

	t.Run("TransferMoneyHandler to return 500 InternalServerError for invalid recipient email", func(t *testing.T) {
		sender := &models.User{
			EmailID:  "sender@example.com",
			Password: "test123",
		}
		senderID, _ := userService.CreateUser(sender)
		_, _ = walletService.CreateWallet(senderID, money.INR)

		invalidRecipientEmail := "invalidemail"
		transferMoney, _ := money.NewMoney(decimal.NewFromFloat(50.0), money.INR)

		transferMoneyPayload := dto.MoneyTransferDTO{
			Amount:         transferMoney,
			RecipientEmail: invalidRecipientEmail,
		}

		reqBody, _ := json.Marshal(transferMoneyPayload)

		url := "/wallet/transfer"
		req, _ := http.NewRequest("PUT", url, bytes.NewReader(reqBody))

		IDToken, _ := authService.AuthenticateUser(sender.EmailID, sender.Password)
		req.Header.Set("id_token", IDToken)

		recorder := httptest.NewRecorder()
		http.HandlerFunc(walletHandlers.TransferMoneyHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

	t.Run("TransferMoneyHandler to return 400 BadRequest for invalid payload", func(t *testing.T) {
		sender := &models.User{
			EmailID:  "sender@example.com",
			Password: "test123",
		}
		senderID, _ := userService.CreateUser(sender)
		_, _ = walletService.CreateWallet(senderID, money.INR)

		IDToken, _ := authService.AuthenticateUser(sender.EmailID, sender.Password)

		invalidPayload := []byte(`{"amount": "100", "recipient_email": "recipient@example.com"}`)

		url := "/wallet/transfer"
		req, _ := http.NewRequest("PUT", url, bytes.NewReader(invalidPayload))
		req.Header.Set("id_token", IDToken)

		recorder := httptest.NewRecorder()
		http.HandlerFunc(walletHandlers.TransferMoneyHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var response dto.Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		assert.Error(t, err, "invalid amount")
	})
}
