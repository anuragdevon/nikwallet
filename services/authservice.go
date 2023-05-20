package services

import (
	"nikwallet/database"

	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(database *gorm.DB) *AuthService {
	return &AuthService{db: database}
}

func (as *AuthService) AuthenticateUser(email string, password string) (string, error) {
	db := database.PostgreSQL{DB: as.db}
	return db.AuthenticateUser(email, password)
}

func (as *AuthService) VerifyToken(tokenString string) (*database.Claims, int, error) {
	db := database.PostgreSQL{DB: as.db}
	return db.VerifyToken(tokenString)
}
