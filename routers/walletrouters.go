package routers

import (
	"net/http"

	"nikwallet/handlers"

	"github.com/gorilla/mux"
)

func NewWalletRouter(handlers *handlers.WalletHandlers) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("", handlers.CreateWalletHandler).Methods(http.MethodPost)
	router.HandleFunc("", handlers.AddMoneyToWalletHandler).Methods(http.MethodPut)
	router.HandleFunc("/withdraw", handlers.WithdrawMoneyFromWalletHandler).Methods(http.MethodPut)
	router.HandleFunc("/transfer", handlers.WithdrawMoneyFromWalletHandler).Methods(http.MethodPut)

	return router
}
