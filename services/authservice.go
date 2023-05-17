package services

import (
	"database/sql"
	"nikwallet/pkg/auth"
)

type AuthService struct {
	db *sql.DB
}

func NewAuthService(database *sql.DB) *AuthService {
	return &AuthService{db: database}
}

func (as *AuthService) AuthenticateUser(email string, password string) (string, error) {
	return auth.AuthenticateUser(email, password)
}

func (as *AuthService) VerifyToken(tokenString string) (*auth.Claims, int, error) {
	return auth.VerifyToken(tokenString)
}
