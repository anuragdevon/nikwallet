package handlers

import (
	"encoding/json"
	"net/http"
	"nikwallet/handlers/dto"
	"nikwallet/repository/money"
	"nikwallet/services"
	"strconv"
)

type WalletHandlers struct {
	walletService *services.WalletService
	authService   *services.AuthService
	userService   *services.UserService
}

func NewWalletHandlers(walletService *services.WalletService, authService *services.AuthService, userService *services.UserService) *WalletHandlers {
	return &WalletHandlers{
		walletService: walletService,
		authService:   authService,
		userService:   userService,
	}
}

func (wh *WalletHandlers) CreateWalletHandler(respWriter http.ResponseWriter, req *http.Request) {
	IDToken := req.Header.Get("id_token")
	_, userID, err := wh.authService.VerifyToken(IDToken)
	if err != nil {
		respWriter.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(respWriter).Encode(dto.Response{Error: err.Error()})
		return
	}

	var payload struct {
		Currency money.Currency `json:"currency"`
	}

	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		http.Error(respWriter, "invalid payload", http.StatusBadRequest)
		return
	}

	wallet, err := wh.walletService.CreateWallet(userID, payload.Currency)
	if err != nil {
		respWriter.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(respWriter).Encode(dto.Response{Error: err.Error()})
		return
	}

	respWriter.WriteHeader(http.StatusCreated)
	json.NewEncoder(respWriter).Encode(wallet)
}

func (wh *WalletHandlers) AddMoneyToWalletHandler(respWriter http.ResponseWriter, req *http.Request) {
	IDToken := req.Header.Get("id_token")
	_, userID, err := wh.authService.VerifyToken(IDToken)
	if err != nil {
		respWriter.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(respWriter).Encode(dto.Response{Error: err.Error()})
		return
	}

	var moneyToAdd money.Money
	if err := json.NewDecoder(req.Body).Decode(&moneyToAdd); err != nil {
		http.Error(respWriter, "invalid amount", http.StatusBadRequest)
		return
	}
	updatedWallet, err := wh.walletService.AddMoneyToWallet(userID, moneyToAdd)
	if err != nil {
		respWriter.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(respWriter).Encode(dto.Response{Error: err.Error()})
		return
	}

	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(&updatedWallet)
}

func (wh *WalletHandlers) WithdrawMoneyFromWalletHandler(respWriter http.ResponseWriter, req *http.Request) {
	IDToken := req.Header.Get("id_token")
	_, userID, err := wh.authService.VerifyToken(IDToken)
	if err != nil {
		respWriter.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(respWriter).Encode(dto.Response{Error: err.Error()})
		return
	}

	var moneyToAdd money.Money
	if err := json.NewDecoder(req.Body).Decode(&moneyToAdd); err != nil {
		http.Error(respWriter, "invalid amount", http.StatusBadRequest)
		return
	}
	var withdrawnMoney money.Money

	if withdrawnMoney, err = wh.walletService.WithdrawMoneyFromWallet(userID, moneyToAdd); err != nil {
		http.Error(respWriter, "insufficient funds", http.StatusBadRequest)
		return
	}

	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(&withdrawnMoney)
}

func (wh *WalletHandlers) TransferMoneyHandler(respWriter http.ResponseWriter, req *http.Request) {
	IDToken := req.Header.Get("id_token")
	_, userID, err := wh.authService.VerifyToken(IDToken)
	if err != nil {
		respWriter.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(respWriter).Encode(dto.Response{Error: err.Error()})
		return
	}

	var transferPayload dto.MoneyTransferDTO

	if err := json.NewDecoder(req.Body).Decode(&transferPayload); err != nil {
		http.Error(respWriter, "invalid payload", http.StatusBadRequest)
		return
	}

	if err := wh.walletService.TransferMoney(userID, transferPayload.RecipientEmail, *transferPayload.Amount); err != nil {
		http.Error(respWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(dto.Response{Message: "money transferred successfully"})
}

func (wh *WalletHandlers) GetWalletHistoryHandler(respWriter http.ResponseWriter, req *http.Request) {
	IDToken := req.Header.Get("id_token")
	_, userID, err := wh.authService.VerifyToken(IDToken)
	if err != nil {
		respWriter.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(respWriter).Encode(dto.Response{Error: err.Error()})
		return
	}

	limitStr := req.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(respWriter, "invalid limit parameter", http.StatusBadRequest)
		return
	}
	ledgerEntries, err := wh.walletService.GetLastNLedgerEntries(userID, limit)
	if err != nil {
		respWriter.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(respWriter).Encode(dto.Response{Error: err.Error()})
		return
	}

	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(ledgerEntries)
}
