package user

import (
	"fmt"
	"nikwallet/pkg/db"
	"testing"
)

func TestCreateUserToCreateAValidUser(t *testing.T) {
	database, err := db.ConnectToDB("testdb")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	user := &User{
		EmailID:  "test2@example.com",
		Password: "test123",
	}
	userID, err := CreateUser(database, user)
	if err != nil {
		t.Errorf("CreateUser() error = %v, want nil", err)
		return
	}

	if userID == 0 {
		t.Errorf("CreateUser() did not set user ID")
	}
}

func TestCreateUserToReturnErrorWithDuplicateEmail(t *testing.T) {
	db, err := db.ConnectToDB("testdb")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	user := &User{
		EmailID:  "anuragkar1@gmail.com",
		Password: "password123",
	}

	_, err = CreateUser(db, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	duplicateUser := &User{
		EmailID:  "anuragkar1@gmail.com",
		Password: "password456",
	}

	_, err = CreateUser(db, duplicateUser)
	if err == nil {
		t.Fatalf("Expected to return err with duplicate email")
	}
}

func TestGetUserByIDToReturnValidUser(t *testing.T) {
	database, err := db.ConnectToDB("testdb")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	user := &User{
		EmailID:  "test4@example.com",
		Password: "test123",
	}
	userID, err := CreateUser(database, user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	fetchedUser, err := GetUserByID(database, userID)
	if err != nil {
		t.Fatalf("GetUserByID() error = %v, want nil", err)
	}

	if fetchedUser.EmailID != user.EmailID {
		t.Errorf("GetUserByID() EmailID = %v, want %v", fetchedUser.EmailID, user.EmailID)
	}
}

func TestGetUserByEmail(t *testing.T) {
	database, err := db.ConnectToDB("testdb")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	// Insert a test user into the database
	userEmail := "testuser99@example.com"
	userPassword := "password123"
	userID, err := CreateUser(database, &User{EmailID: userEmail, Password: userPassword})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Test getting the user by email
	user, err := GetUserByEmail(database, userEmail)
	if err != nil {
		t.Fatalf("failed to get user by email: %v", err)
	}

	if user.ID != userID {
		t.Errorf("GetUserByEmail() returned wrong user ID, got %d, want %d", user.ID, userID)
	}

	if user.EmailID != userEmail {
		t.Errorf("GetUserByEmail() returned wrong email, got %s, want %s", user.EmailID, userEmail)
	}
}

func TestGetUserByEmailNotFound(t *testing.T) {
	database, err := db.ConnectToDB("testdb")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	// Test getting a user that doesn't exist
	userEmail := "nonexistent@example.com"
	_, err = GetUserByEmail(database, userEmail)
	if err == nil {
		t.Fatalf("expected GetUserByEmail() to return an error, but got nil")
	}

	expectedErrorMessage := fmt.Sprintf("user with email %s not found", userEmail)
	if err.Error() != expectedErrorMessage {
		t.Errorf("GetUserByEmail() returned wrong error message, got %s, want %s", err.Error(), expectedErrorMessage)
	}
}
