package handlers

import (
	"encoding/json"
	"net/http"
	"nikwallet/pkg/money"
	"nikwallet/services"
	"strconv"

	"github.com/gorilla/mux"
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

type Response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (wh *WalletHandlers) CreateWalletHandler(w http.ResponseWriter, r *http.Request) {
	IDToken := r.Header.Get("id_token")
	_, userID, err := wh.authService.VerifyToken(IDToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{Error: err.Error()})
		return
	}

	wallet, err := wh.walletService.CreateWallet(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(wallet)
}

func (wh *WalletHandlers) AddMoneyToWalletHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	IDToken := r.Header.Get("id_token")
	_, _, err := wh.authService.VerifyToken(IDToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{Error: err.Error()})
		return
	}
	walletID, err := strconv.Atoi(vars["wallet_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Error: "invalid wallet ID"})
		return
	}

	var m money.Money
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Error: "invalid amount"})
		return
	}

	if err := wh.walletService.AddMoneyToWallet(walletID, m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Message: "money added to wallet successfully"})
}

func (wh *WalletHandlers) WithdrawMoneyFromWalletHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	IDToken := r.Header.Get("id_token")
	_, _, err := wh.authService.VerifyToken(IDToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{Error: err.Error()})
		return
	}

	walletID, err := strconv.Atoi(vars["wallet_id"])
	if err != nil {
		http.Error(w, "invalid wallet ID", http.StatusBadRequest)
		return
	}

	var m money.Money
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	if _, err = wh.walletService.WithdrawMoneyFromWallet(walletID, m); err != nil {
		http.Error(w, "insufficient funds", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}