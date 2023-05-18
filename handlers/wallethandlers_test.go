package handlers_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
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

		url := "/wallet/add"
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

	t.Run("WithdrawMoneyFromWalletHandler to return status 400 bad request for InvalidAmount", func(t *testing.T) {
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

}
