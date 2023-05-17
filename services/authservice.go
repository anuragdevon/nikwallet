package services

import (
	"database/sql"
	"nikwallet/database"
)

type AuthService struct {
	db *sql.DB
}

func NewAuthService(database *sql.DB) *AuthService {
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
