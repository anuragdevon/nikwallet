package routers

import (
	"net/http"

	"nikwallet/handlers"
	"nikwallet/services"

	"github.com/gorilla/mux"
)

func NewUserRouter(userService *services.UserService, authService *services.AuthService) *mux.Router {
	router := mux.NewRouter()
	handlers := handlers.NewUserHandlers(userService, authService)

	router.HandleFunc("/signup", handlers.SignupHandler).Methods(http.MethodPost)
	router.HandleFunc("/signin", handlers.SigninHandler).Methods(http.MethodPost)

	return router
}
