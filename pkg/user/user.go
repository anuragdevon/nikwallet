package user

import (
	"fmt"
	"nikwallet/pkg/db"
)

type User struct {
	ID       int
	EmailID  string
	Password string
}

func CreateUser(user *User) (int, error) {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	err := db.DB.QueryRow(query, user.EmailID, user.Password).Scan(&user.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return user.ID, nil
}

func GetUserByID(id int) (*User, error) {
	query := `SELECT id, email, password FROM users WHERE id=$1`
	user := &User{}

	err := db.DB.QueryRow(query, id).Scan(&user.ID, &user.EmailID, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("user with id %d not found", id)
	}
	return user, nil
}

func GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, email, password FROM users WHERE email=$1`
	user := &User{}

	err := db.DB.QueryRow(query, email).Scan(&user.ID, &user.EmailID, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("user with email %s not found", email)
	}

	return user, nil
}
