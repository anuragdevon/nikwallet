package handlers_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
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

	newUser := &database.User{
		EmailID:  "testw5111@example.com",
		Password: "password",
	}

	t.Run("TestCreateWalletHandlerToSuccessfullyCreateWallet", func(t *testing.T) {
		_, err = userService.CreateUser(newUser)
		assert.NoError(t, err)

		token, err := authService.AuthenticateUser(newUser.EmailID, newUser.Password)
		assert.Nil(t, err)
		assert.NotNil(t, token)

		req, err := http.NewRequest("POST", "/create", nil)
		req.Header.Set("id_token", token)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(walletHandlers.CreateWalletHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)

		var walletID int
		err = json.NewDecoder(recorder.Body).Decode(&walletID)
		assert.NoError(t, err)
	})

	t.Run("TestAddMoneyToWalletHandlerToSuccessfullyAddMoney", func(t *testing.T) {
		userID, err := userService.CreateUser(newUser)
		assert.NoError(t, err)

		walletID, err := walletService.CreateWallet(userID)
		assert.NoError(t, err)

		IDToken, _ := authService.AuthenticateUser(newUser.EmailID, newUser.Password)
		assert.NotNil(t, IDToken)

		addMoneyRequest := money.Money{Amount: 50, Currency: "INR"}
		reqBody, err := json.Marshal(addMoneyRequest)
		assert.NoError(t, err)

		url := "/add/" + strconv.Itoa(walletID)
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

	t.Run("TestWithdrawMoneyFromWalletHandlerToSuccessfullyWithdrawMoney", func(t *testing.T) {
		userService := services.NewUserService(db.DB)
		authService := services.NewAuthService(db.DB)
		walletService := services.NewWalletService(db.DB)

		walletHandlers := handlers.NewWalletHandlers(walletService, authService, userService)

		newUser := &database.User{
			EmailID:  "testw5111@example.com",
			Password: "password",
		}

		userID, err := userService.CreateUser(newUser)
		assert.NoError(t, err)

		walletID, err := walletService.CreateWallet(userID)
		assert.NoError(t, err)

		IDToken, _ := authService.AuthenticateUser(newUser.EmailID, newUser.Password)
		assert.NotNil(t, IDToken)

		addMoneyRequest := money.Money{Amount: 50, Currency: "INR"}

		err = walletService.AddMoneyToWallet(walletID, addMoneyRequest)
		assert.NoError(t, err)

		withdrawMoneyRequest := money.Money{Amount: 50, Currency: "INR"}
		reqBody, err := json.Marshal(withdrawMoneyRequest)
		assert.NoError(t, err)

		url := "/wallet/" + strconv.Itoa(walletID) + "/withdraw"
		req, err := http.NewRequest("PUT", url, bytes.NewReader(reqBody))
		req.Header.Set("id_token", IDToken)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(walletHandlers.WithdrawMoneyFromWalletHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response handlers.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, "money withdrawn from wallet successfully", response.Message)
	})

}
