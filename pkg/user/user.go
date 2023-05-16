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

func CreateUser(database *db.DB, user *User) (int, error) {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	err := database.QueryRow(query, user.EmailID, user.Password).Scan(&user.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return user.ID, nil
}

func GetUserByID(database *db.DB, id int) (*User, error) {
	query := `SELECT id, email, password FROM users WHERE id=$1`
	user := &User{}

	err := database.QueryRow(query, id).Scan(&user.ID, &user.EmailID, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("user with id %d not found", id)
	}
	return user, nil
}
