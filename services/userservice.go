package services

import (
	"fmt"
	"nikwallet/repository"
	"nikwallet/repository/models"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (us *UserService) CreateUser(newUser *models.User) (int, error) {
	db := repository.PostgreSQL{DB: us.db}

	existingUser, _ := db.GetUserByEmail(newUser.EmailID)

	if existingUser != nil {
		return 0, fmt.Errorf("user already exists")
	}

	return db.CreateUser(newUser)
}

func (us *UserService) GetUserByID(id int) (*models.User, error) {
	db := repository.PostgreSQL{DB: us.db}
	return db.GetUserByID(id)
}

func (us *UserService) GetUserByEmail(email string) (*models.User, error) {
	db := repository.PostgreSQL{DB: us.db}
	return db.GetUserByEmail(email)
}
