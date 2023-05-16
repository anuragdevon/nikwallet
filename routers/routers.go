package routers

import (
	"net/http"

	"nikwallet/services"

	"github.com/gorilla/mux"
)

func NewRouter(userService *services.UserService, authService *services.AuthService, walletService *services.WalletService) *mux.Router {
	router := mux.NewRouter()

	userRouter := NewUserRouter(userService, authService)
	router.PathPrefix("/user").Handler(http.StripPrefix("/user", userRouter))

	walletRouter := NewWalletRouter(walletService, authService)
	router.PathPrefix("/wallet").Handler(http.StripPrefix("/wallet", walletRouter))

	return router
}
