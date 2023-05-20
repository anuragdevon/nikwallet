package handlers

import (
	"encoding/json"
	"net/http"

	"nikwallet/repository/models"
	"nikwallet/services"
)

type UserHandlers struct {
	userService *services.UserService
	authService *services.AuthService
}

func NewUserHandlers(userService *services.UserService, authService *services.AuthService) *UserHandlers {
	return &UserHandlers{
		userService: userService,
		authService: authService,
	}
}

func (uh *UserHandlers) SignupHandler(respWriter http.ResponseWriter, req *http.Request) {
	var userData models.User

	if err := json.NewDecoder(req.Body).Decode(&userData); err != nil {
		http.Error(respWriter, err.Error(), http.StatusBadRequest)
		return
	}
	createdUserID, err := uh.userService.CreateUser(&userData)
	if err != nil {
		http.Error(respWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenString, _ := uh.authService.AuthenticateUser(userData.EmailID, userData.Password)

	respWriter.WriteHeader(http.StatusCreated)
	json.NewEncoder(respWriter).Encode(map[string]interface{}{
		"user_id": createdUserID,
		"token":   tokenString,
	})
}

func (uh *UserHandlers) SigninHandler(respWriter http.ResponseWriter, req *http.Request) {
	var userData models.User

	if err := json.NewDecoder(req.Body).Decode(&userData); err != nil {
		http.Error(respWriter, err.Error(), http.StatusBadRequest)
		return
	}

	tokenString, err := uh.authService.AuthenticateUser(userData.EmailID, userData.Password)
	if err != nil {
		http.Error(respWriter, err.Error(), http.StatusUnauthorized)
		return
	}

	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(map[string]string{
		"token": tokenString,
	})
}
