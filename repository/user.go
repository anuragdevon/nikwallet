package repository

import (
	"fmt"
	"nikwallet/repository/models"
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
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (db *PostgreSQL) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := db.DB.Where("email_id = ?", email).First(user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
