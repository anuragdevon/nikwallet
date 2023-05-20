package services

import (
	"nikwallet/repository"

	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (as *AuthService) AuthenticateUser(email string, password string) (string, error) {
	db := repository.PostgreSQL{DB: as.db}
	return db.AuthenticateUser(email, password)
}

func (as *AuthService) VerifyToken(tokenString string) (*repository.Claims, int, error) {
	db := repository.PostgreSQL{DB: as.db}
	return db.VerifyToken(tokenString)
}
