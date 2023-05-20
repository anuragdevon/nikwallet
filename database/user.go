package database

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	EmailID  string `gorm:"unique"`
	Password string
}

func (db *PostgreSQL) CreateUser(newUser *User) (int, error) {
	err := db.DB.Create(newUser).Error
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return int(newUser.ID), nil
}

func (db *PostgreSQL) GetUserByID(id int) (*User, error) {
	user := &User{}
	err := db.DB.First(user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (db *PostgreSQL) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := db.DB.Where("email_id = ?", email).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
