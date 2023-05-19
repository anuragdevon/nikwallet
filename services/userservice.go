package services

import (
	"nikwallet/database"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
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
