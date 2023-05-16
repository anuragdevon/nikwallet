package services

import (
	"nikwallet/pkg/auth"
	"nikwallet/pkg/db"
)

type AuthService struct {
	db *db.DB
}

func NewAuthService(database *db.DB) *AuthService {
	return &AuthService{db: database}
}

func (as *AuthService) AuthenticateUser(email string, password string) (string, error) {
	return auth.AuthenticateUser(as.db, email, password)
}

func (as *AuthService) VerifyToken(tokenString string) (*auth.Claims, int, error) {
	return auth.VerifyToken(tokenString)
}
