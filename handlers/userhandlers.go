package handlers

import (
	"encoding/json"
	"net/http"

	"nikwallet/handlers/dto"
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
	var userData dto.UserSignupRequestDTO

	if err := json.NewDecoder(req.Body).Decode(&userData); err != nil {
		http.Error(respWriter, err.Error(), http.StatusBadRequest)
		return
	}

	user := models.User{
		EmailID:  userData.Email,
		Password: userData.Password,
	}

	createdUserID, err := uh.userService.CreateUser(&user)
	if err != nil {
		http.Error(respWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenString, _ := uh.authService.AuthenticateUser(user.EmailID, user.Password)

	responseDTO := dto.UserSignupResponseDTO{
		UserID:  createdUserID,
		IDToken: tokenString,
	}

	respWriter.WriteHeader(http.StatusCreated)
	json.NewEncoder(respWriter).Encode(responseDTO)
}

func (uh *UserHandlers) SigninHandler(respWriter http.ResponseWriter, req *http.Request) {
	var userData dto.UserSigninRequestDTO

	if err := json.NewDecoder(req.Body).Decode(&userData); err != nil {
		http.Error(respWriter, err.Error(), http.StatusBadRequest)
		return
	}

	tokenString, err := uh.authService.AuthenticateUser(userData.Email, userData.Password)
	if err != nil {
		http.Error(respWriter, err.Error(), http.StatusUnauthorized)
		return
	}

	responseDTO := dto.UserSigninResponseDTO{
		IDToken: tokenString,
	}

	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(responseDTO)
}
