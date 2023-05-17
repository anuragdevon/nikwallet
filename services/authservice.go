package services

import (
	"nikwallet/pkg/auth"
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (as *AuthService) AuthenticateUser(email string, password string) (string, error) {
	return auth.AuthenticateUser(email, password)
}

func (as *AuthService) VerifyToken(tokenString string) (*auth.Claims, int, error) {
	return auth.VerifyToken(tokenString)
}
