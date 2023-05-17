package services

import (
	"database/sql"
	"nikwallet/pkg/user"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (us *UserService) CreateUser(newUser *user.User) (int, error) {
	return user.CreateUser(newUser)
}

func (us *UserService) GetUserByID(id int) (*user.User, error) {
	return user.GetUserByID(id)
}

func (us *UserService) GetUserByEmail(email string) (*user.User, error) {
	return user.GetUserByEmail(email)
}
