package handlers_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"nikwallet/database"
	"nikwallet/database/money"
	"nikwallet/handlers"
	"nikwallet/services"
)

func TestWalletHandlers(t *testing.T) {
	db := &database.PostgreSQL{}
	err := db.Connect("testdb")
	if err != nil {
		log.Panic("failed to connect to database:", err)
	}
	defer db.Close()

	userService := services.NewUserService(db.DB)
	authService := services.NewAuthService(db.DB)
	walletService := services.NewWalletService(db.DB)

	walletHandlers := handlers.NewWalletHandlers(walletService, authService, userService)

	t.Run("CreateWalletHandler to return 201 StatusCreated for successfully create wallet", func(t *testing.T) {
		newUser := &database.User{
			EmailID:  "testw5111@example.com",
			Password: "password",
		}

		_, err = userService.CreateUser(newUser)
		assert.NoError(t, err)

		token, err := authService.AuthenticateUser(newUser.EmailID, newUser.Password)
		assert.Nil(t, err)
		assert.NotNil(t, token)

		req, err := http.NewRequest("POST", "/wallet/create", nil)
		req.Header.Set("id_token", token)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(walletHandlers.CreateWalletHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)

		var wallet *database.Wallet
		err = json.NewDecoder(recorder.Body).Decode(&wallet)
		assert.NoError(t, err)
	})

	t.Run("AddMoneyToWalletHandler to return 200 StatusOk for successfull add money to user's wallet", func(t *testing.T) {
		newUser := &database.User{
			EmailID:  "testw5112@example.com",
			Password: "password",
		}

		userID, err := userService.CreateUser(newUser)
		assert.NoError(t, err)

		_, err = walletService.CreateWallet(userID)
		assert.NoError(t, err)

		IDToken, _ := authService.AuthenticateUser(newUser.EmailID, newUser.Password)
		assert.NotNil(t, IDToken)

		addMoneyRequest := money.Money{Amount: 50, Currency: "INR"}
		reqBody, err := json.Marshal(addMoneyRequest)
		assert.NoError(t, err)

		url := "/wallet"
		req, err := http.NewRequest("PUT", url, bytes.NewReader(reqBody))
		req.Header.Set("id_token", IDToken)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(walletHandlers.AddMoneyToWalletHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response handlers.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, "money added to wallet successfully", response.Message)
	})

	t.Run("AddMoneyToWalletHandler to return status 400 bad request for InvalidAmount", func(t *testing.T) {
		newUser := &database.User{
			EmailID:  "testw599@example.com",
			Password: "password",
		}

		userID, err := userService.CreateUser(newUser)
		assert.NoError(t, err)

		_, err = walletService.CreateWallet(userID)
		assert.NoError(t, err)

		IDToken, _ := authService.AuthenticateUser(newUser.EmailID, newUser.Password)
		assert.NotNil(t, IDToken)

		invalidWithdrawMoneyRequest := map[string]interface{}{
			"amount":   "notanumber",
			"currency": "INR",
		}

		reqBody, err := json.Marshal(invalidWithdrawMoneyRequest)
		assert.NoError(t, err)

		url := "/wallet/add"
		req, err := http.NewRequest("PUT", url, bytes.NewReader(reqBody))
		req.Header.Set("id_token", IDToken)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(walletHandlers.AddMoneyToWalletHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var response handlers.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		assert.Error(t, err, "invalid amount")
	})

	t.Run("WithdrawMoneyFromWalletHandler to return 200 StatusOk for successfull withdraw money from user's wallet", func(t *testing.T) {
		newUser := &database.User{
			EmailID:  "testw5113@example.com",
			Password: "password",
		}

		userID, err := userService.CreateUser(newUser)
		assert.NoError(t, err)

		_, err = walletService.CreateWallet(userID)
		assert.NoError(t, err)

		IDToken, _ := authService.AuthenticateUser(newUser.EmailID, newUser.Password)
		assert.NotNil(t, IDToken)

		addMoneyRequest := money.Money{Amount: 50, Currency: "INR"}

		err = walletService.AddMoneyToWallet(userID, addMoneyRequest)
		assert.NoError(t, err)

		withdrawMoneyRequest := money.Money{Amount: 50, Currency: "INR"}
		reqBody, err := json.Marshal(withdrawMoneyRequest)
		assert.NoError(t, err)

		url := "wallet/withdraw"
		req, err := http.NewRequest("PUT", url, bytes.NewReader(reqBody))
		req.Header.Set("id_token", IDToken)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(walletHandlers.WithdrawMoneyFromWalletHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var wallet *database.Wallet
		err = json.NewDecoder(recorder.Body).Decode(&wallet)
		assert.NoError(t, err)
	})

	t.Run("WithdrawMoneyFromWalletHandler to return status 400 bad request for InsufficientFunds", func(t *testing.T) {
		newUser := &database.User{
			EmailID:  "testw5114@example.com",
			Password: "password",
		}

		userID, err := userService.CreateUser(newUser)
		assert.NoError(t, err)

		_, err = walletService.CreateWallet(userID)
		assert.NoError(t, err)

		IDToken, _ := authService.AuthenticateUser(newUser.EmailID, newUser.Password)
		assert.NotNil(t, IDToken)

		addMoneyRequest := money.Money{Amount: 40, Currency: "INR"}
		err = walletService.AddMoneyToWallet(userID, addMoneyRequest)
		assert.NoError(t, err)

		withdrawMoneyRequest := money.Money{Amount: 50, Currency: "INR"}
		reqBody, err := json.Marshal(withdrawMoneyRequest)
		assert.NoError(t, err)

		url := "/wallet/withdraw"
		req, err := http.NewRequest("PUT", url, bytes.NewReader(reqBody))
		req.Header.Set("id_token", IDToken)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(walletHandlers.WithdrawMoneyFromWalletHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var response handlers.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		assert.Error(t, err, "insufficient funds")
	})

	t.Run("WithdrawMoneyFromWalletHandler to return status 400 BadRequest for InvalidAmount", func(t *testing.T) {
		newUser := &database.User{
			EmailID:  "testw5115@example.com",
			Password: "password",
		}

		userID, err := userService.CreateUser(newUser)
		assert.NoError(t, err)

		_, err = walletService.CreateWallet(userID)
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

		var response handlers.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		assert.Error(t, err, "invalid amount")
	})

	t.Run("TransferMoneyHandler to return 200 StatusOk for successful transfer of money from sender to reciever", func(t *testing.T) {
		sender := &database.User{
			EmailID:  "sender@example.com",
			Password: "test123",
		}
		senderID, _ := userService.CreateUser(sender)
		_, _ = walletService.CreateWallet(senderID)

		recipient := &database.User{
			EmailID:  "recipient@example.com",
			Password: "test123",
		}
		recipientID, _ := userService.CreateUser(recipient)
		_, _ = walletService.CreateWallet(recipientID)

		initialMoney, _ := money.NewMoney(100, "INR")
		_ = walletService.AddMoneyToWallet(senderID, *initialMoney)

		IDToken, _ := authService.AuthenticateUser(sender.EmailID, sender.Password)

		transferMoney, _ := money.NewMoney(50, "INR")

		transferMoneyPayload := map[string]interface{}{
			"amount":          transferMoney,
			"recipient_email": recipient.EmailID,
		}

		reqBody, _ := json.Marshal(transferMoneyPayload)

		url := "/wallet/transfer"
		req, _ := http.NewRequest("PUT", url, bytes.NewReader(reqBody))
		req.Header.Set("id_token", IDToken)

		recorder := httptest.NewRecorder()
		http.HandlerFunc(walletHandlers.TransferMoneyHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		var response handlers.Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "money transferred successfully", response.Message)

		expectedSenderMoney, _ := money.NewMoney(50, "INR")
		senderWallet, _ := walletService.GetWalletByUserID(senderID)
		if !reflect.DeepEqual(&senderWallet.Money, expectedSenderMoney) {
			t.Errorf("TransferMoneyHandler() sender balance got = %v, want = %v", senderWallet.Money, expectedSenderMoney)
		}

		expectedRecipientMoney, _ := money.NewMoney(50, "INR")
		recipientWallet, _ := walletService.GetWalletByUserID(recipientID)
		if !reflect.DeepEqual(&recipientWallet.Money, expectedRecipientMoney) {
			t.Errorf("TransferMoneyHandler() recipient balance got = %v, want = %v", recipientWallet.Money, expectedRecipientMoney)
		}
	})

	t.Run("TransferMoneyHandler to return 500 InternalServerError for invalid recipient email", func(t *testing.T) {
		sender := &database.User{
			EmailID:  "sender@example.com",
			Password: "test123",
		}
		senderID, _ := userService.CreateUser(sender)
		_, _ = walletService.CreateWallet(senderID)

		invalidRecipientEmail := "invalidemail"
		transferMoney, _ := money.NewMoney(50, "INR")

		transferMoneyPayload := map[string]interface{}{
			"amount":          transferMoney,
			"recipient_email": invalidRecipientEmail,
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
		sender := &database.User{
			EmailID:  "sender@example.com",
			Password: "test123",
		}
		senderID, _ := userService.CreateUser(sender)
		_, _ = walletService.CreateWallet(senderID)

		IDToken, _ := authService.AuthenticateUser(sender.EmailID, sender.Password)

		invalidPayload := []byte(`{"amount": "100", "recipient_email": "recipient@example.com"}`)

		url := "/wallet/transfer"
		req, _ := http.NewRequest("PUT", url, bytes.NewReader(invalidPayload))
		req.Header.Set("id_token", IDToken)

		recorder := httptest.NewRecorder()
		http.HandlerFunc(walletHandlers.TransferMoneyHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var response handlers.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		assert.Error(t, err, "invalid amount")
	})
}
