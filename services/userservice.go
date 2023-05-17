package services

import (
	"database/sql"
	"nikwallet/database"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (us *UserService) CreateUser(newUser *database.User) (int, error) {
	db := database.PostgreSQL{DB: us.db}
	return db.CreateUser(newUser)
}

func (us *UserService) GetUserByID(id int) (*database.User, error) {
	db := database.PostgreSQL{DB: us.db}
	return db.GetUserByID(id)
}

func (us *UserService) GetUserByEmail(email string) (*database.User, error) {
	db := database.PostgreSQL{DB: us.db}
	return db.GetUserByEmail(email)
}
