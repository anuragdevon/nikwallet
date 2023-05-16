package routers

import (
	"nikwallet/services"

	"github.com/gorilla/mux"
)

func NewWalletRouter(walletService *services.WalletService, authService *services.AuthService) *mux.Router {
	router := mux.NewRouter()

	return router
}
