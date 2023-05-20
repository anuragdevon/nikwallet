package database

import (
	"fmt"
)

type User struct {
	ID       int
	EmailID  string `json:"email_id"`
	Password string `json:"password"`
}

func (db *PostgreSQL) CreateUser(newUser *User) (int, error) {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	err := db.DB.QueryRow(query, newUser.EmailID, newUser.Password).Scan(&newUser.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return newUser.ID, nil
}

func (db *PostgreSQL) GetUserByID(id int) (*User, error) {
	query := `SELECT id, email, password FROM users WHERE id=$1`
	user := &User{}

	err := db.DB.QueryRow(query, id).Scan(&user.ID, &user.EmailID, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("user with id %d not found", id)
	}
	return user, nil
}

func (db *PostgreSQL) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, email, password FROM users WHERE email=$1`
	user := &User{}

	err := db.DB.QueryRow(query, email).Scan(&user.ID, &user.EmailID, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("user with email %s not found", email)
	}

	return user, nil
}
