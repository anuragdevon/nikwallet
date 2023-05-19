package database

import (
	"fmt"
	"testing"
)

func TestUser(t *testing.T) {
	t.Run("CreateUser method to successfully create user with valid data", func(t *testing.T) {
		user := &User{
			EmailID:  "test@example.com",
			Password: "test123",
		}
		userID, err := db.CreateUser(user)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		if userID == 0 {
			t.Errorf("CreateUser() did not set user ID")
		}
	})

	t.Run("CreateUser method to return error for duplicate emailID", func(t *testing.T) {
		user := &User{
			EmailID:  "anuragkar1@gmail.com",
			Password: "password123",
		}

		_, err := db.CreateUser(user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		duplicateUser := &User{
			EmailID:  "anuragkar1@gmail.com",
			Password: "password456",
		}

		_, err = db.CreateUser(duplicateUser)
		if err == nil {
			t.Fatalf("Expected to return err with duplicate email")
		}
	})

	t.Run("GetUserByID method to return valid user for valid userID", func(t *testing.T) {
		user := &User{
			EmailID:  "test4@example.com",
			Password: "test123",
		}
		userID, err := db.CreateUser(user)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		fetchedUser, err := db.GetUserByID(userID)
		if err != nil {
			t.Fatalf("GetUserByID() error = %v, want nil", err)
		}

		if fetchedUser.EmailID != user.EmailID {
			t.Errorf("GetUserByID() EmailID = %v, want %v", fetchedUser.EmailID, user.EmailID)
		}
	})

	t.Run("GetUserByEmail method to return valid user for valid emailID", func(t *testing.T) {
		userEmail := "testuser99@example.com"
		userPassword := "password123"
		userID, err := db.CreateUser(&User{EmailID: userEmail, Password: userPassword})
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}

		user, err := db.GetUserByEmail(userEmail)
		if err != nil {
			t.Fatalf("failed to get user by email: %v", err)
		}

		if user.ID != userID {
			t.Errorf("GetUserByEmail() returned wrong user ID, got %d, want %d", user.ID, userID)
		}

		if user.EmailID != userEmail {
			t.Errorf("GetUserByEmail() returned wrong email, got %s, want %s", user.EmailID, userEmail)
		}
	})

	t.Run("GetUserByEmail method to return error for invalid emailID", func(t *testing.T) {
		userEmail := "nonexistent@example.com"
		_, err := db.GetUserByEmail(userEmail)
		if err == nil {
			t.Fatalf("expected GetUserByEmail() to return an error, but got nil")
		}

		expectedErrorMessage := fmt.Sprintf("user with email %s not found", userEmail)
		if err.Error() != expectedErrorMessage {
			t.Fatalf("GetUserByEmail() error = %v, want %v", err.Error(), expectedErrorMessage)
		}
	})
}
