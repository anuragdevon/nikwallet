package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"nikwallet/handlers"
	"nikwallet/pkg/auth"
	"nikwallet/pkg/db"
	"nikwallet/pkg/money"
	"nikwallet/pkg/user"
	"nikwallet/services"
)

func TestCreateWalletHandlerToSuccessfullyCreateWallet(t *testing.T) {
	database, err := db.ConnectToDB("testdb")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()
	userService := services.NewUserService(database)
	authService := services.NewAuthService(database)
	walletService := services.NewWalletService(database)

	walletHandlers := handlers.NewWalletHandlers(walletService, authService, userService)

	newUser := &user.User{
		EmailID:  "testw5111@example.com",
		Password: "password",
	}

	_, err = userService.CreateUser(newUser)
	assert.NoError(t, err)

	token, err := auth.AuthenticateUser(database, newUser.EmailID, newUser.Password)
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

}

func TestAddMoneyToWalletHandlerToSuccessfullyAddMoney(t *testing.T) {
	database, err := db.ConnectToDB("testdb")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()
	userService := services.NewUserService(database)
	authService := services.NewAuthService(database)
	walletService := services.NewWalletService(database)

	walletHandlers := handlers.NewWalletHandlers(walletService, authService, userService)

	newUser := &user.User{
		EmailID:  "testw5111@example.com",
		Password: "password",
	}

	userID, err := userService.CreateUser(newUser)
	assert.NoError(t, err)

	walletID, err := walletService.CreateWallet(userID)
	assert.NoError(t, err)

	IDToken, _ := auth.AuthenticateUser(database, newUser.EmailID, newUser.Password)
	assert.NotNil(t, IDToken)

	addMoneyRequest := money.Money{Amount: 50, Currency: "INR"}
	reqBody, err := json.Marshal(addMoneyRequest)
	assert.NoError(t, err)

	url := "/add/" + strconv.Itoa(walletID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	req.Header.Set("id_token", IDToken)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()

	http.HandlerFunc(walletHandlers.AddMoneyToWalletHandler).ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response handlers.Response
	err = json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, "money added to wallet successfully", response.Message)
}

func TestWithdrawMoneyFromWalletHandlerToSuccessfullyWithdrawMoney(t *testing.T) {
	database, err := db.ConnectToDB("testdb")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()
	userService := services.NewUserService(database)
	authService := services.NewAuthService(database)
	walletService := services.NewWalletService(database)

	walletHandlers := handlers.NewWalletHandlers(walletService, authService, userService)

	newUser := &user.User{
		EmailID:  "testw5111@example.com",
		Password: "password",
	}

	userID, err := userService.CreateUser(newUser)
	assert.NoError(t, err)

	walletID, err := walletService.CreateWallet(userID)
	assert.NoError(t, err)

	IDToken, _ := auth.AuthenticateUser(database, newUser.EmailID, newUser.Password)
	assert.NotNil(t, IDToken)

	addMoneyRequest := money.Money{Amount: 50, Currency: "INR"}

	err = walletService.AddMoneyToWallet(walletID, addMoneyRequest)
	assert.NoError(t, err)

	withdrawMoneyRequest := money.Money{Amount: 50, Currency: "INR"}
	reqBody, err := json.Marshal(withdrawMoneyRequest)
	assert.NoError(t, err)

	url := "/wallets/" + strconv.Itoa(walletID) + "/withdraw"
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	req.Header.Set("id_token", IDToken)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()

	http.HandlerFunc(walletHandlers.WithdrawMoneyFromWalletHandler).ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response handlers.Response
	err = json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, "money withdrawn from wallet successfully", response.Message)
}