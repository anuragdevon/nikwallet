package routers

import (
	"net/http"

	"nikwallet/handlers"

	"github.com/gorilla/mux"
)

func NewRouter(userHandlers *handlers.UserHandlers, walletHandlers *handlers.WalletHandlers) *mux.Router {
	router := mux.NewRouter()

	userRouter := NewUserRouter(userHandlers)
	router.PathPrefix("/user").Handler(http.StripPrefix("/user", userRouter))

	walletRouter := NewWalletRouter(walletHandlers)
	router.PathPrefix("/wallet").Handler(http.StripPrefix("/wallet", walletRouter))

	return router
}
