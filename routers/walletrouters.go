package routers

import (
	"net/http"

	"nikwallet/handlers"

	"github.com/gorilla/mux"
)

func NewWalletRouter(handlers *handlers.WalletHandlers) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/create", handlers.CreateWalletHandler).Methods(http.MethodPost)
	router.HandleFunc("/add/{id}", handlers.AddMoneyToWalletHandler).Methods(http.MethodPut)
	router.HandleFunc("/withdraw/{id}", handlers.WithdrawMoneyFromWalletHandler).Methods(http.MethodPut)

	return router
}
