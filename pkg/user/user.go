package db

import (
	"fmt"
	"nikwallet/pkg/db"
)

type User struct {
	ID       int
	EmailID  string
	Password string
}

func CreateUser(database *db.DB, user *User) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	err := database.QueryRow(query, user.EmailID, user.Password).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}
