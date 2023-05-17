package handlers

import (
	"encoding/json"
	"net/http"

	"nikwallet/database"
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

func (uh *UserHandlers) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var u database.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdUserID, err := uh.userService.CreateUser(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenString, _ := uh.authService.AuthenticateUser(u.EmailID, u.Password)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": createdUserID,
		"token":   tokenString,
	})
}

func (uh *UserHandlers) SigninHandler(w http.ResponseWriter, r *http.Request) {
	var u database.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenString, err := uh.authService.AuthenticateUser(u.EmailID, u.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}
