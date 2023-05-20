package repository

import (
	"fmt"
	"nikwallet/repository/models"

	"gorm.io/gorm"
)

func (db *PostgreSQL) CreateUser(newUser *models.User) (int, error) {
	err := db.DB.Create(newUser).Error
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return int(newUser.ID), nil
}

func (db *PostgreSQL) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	err := db.DB.First(user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (db *PostgreSQL) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := db.DB.Where("email_id = ?", email).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
