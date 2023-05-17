package services

import (
	"nikwallet/pkg/user"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
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
