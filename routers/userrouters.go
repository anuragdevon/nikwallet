package routers

import (
	"net/http"

	"nikwallet/handlers"

	"github.com/gorilla/mux"
)

func NewUserRouter(handlers *handlers.UserHandlers) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/signup", handlers.SignupHandler).Methods(http.MethodPost)
	router.HandleFunc("/signin", handlers.SigninHandler).Methods(http.MethodPost)

	return router
}
