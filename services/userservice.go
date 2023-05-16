package services

import (
	"nikwallet/pkg/db"
	"nikwallet/pkg/user"
)

type UserService struct {
	db *db.DB
}

func NewUserService(database *db.DB) *UserService {
	return &UserService{db: database}
}

func (us *UserService) CreateUser(newUser *user.User) (int, error) {
	return user.CreateUser(us.db, newUser)
}

func (us *UserService) GetUserByID(id int) (*user.User, error) {
	return user.GetUserByID(us.db, id)
}

func (us *UserService) GetUserByEmail(email string) (*user.User, error) {
	return user.GetUserByEmail(us.db, email)
}
