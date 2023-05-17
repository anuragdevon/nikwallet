package user

import (
	"database/sql"
	"fmt"
	"log"
	"nikwallet/pkg/db"
	"os"
	"testing"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	if err := db.ConnectToDB("testdb"); err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}
	testDB = db.DB

	code := m.Run()

	testDB.Close()

	os.Exit(code)
}

func TestCreateUserToCreateAValidUser(t *testing.T) {
	user := &User{
		EmailID:  "test2@example.com",
		Password: "test123",
	}
	userID, err := CreateUser(user)
	if err != nil {
		t.Errorf("CreateUser() error = %v, want nil", err)
		return
	}

	if userID == 0 {
		t.Errorf("CreateUser() did not set user ID")
	}
}

func TestCreateUserToReturnErrorWithDuplicateEmail(t *testing.T) {
	user := &User{
		EmailID:  "anuragkar1@gmail.com",
		Password: "password123",
	}

	_, err := CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	duplicateUser := &User{
		EmailID:  "anuragkar1@gmail.com",
		Password: "password456",
	}

	_, err = CreateUser(duplicateUser)
	if err == nil {
		t.Fatalf("Expected to return err with duplicate email")
	}
}

func TestGetUserByIDToReturnValidUser(t *testing.T) {
	user := &User{
		EmailID:  "test4@example.com",
		Password: "test123",
	}
	userID, err := CreateUser(user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	fetchedUser, err := GetUserByID(userID)
	if err != nil {
		t.Fatalf("GetUserByID() error = %v, want nil", err)
	}

	if fetchedUser.EmailID != user.EmailID {
		t.Errorf("GetUserByID() EmailID = %v, want %v", fetchedUser.EmailID, user.EmailID)
	}
}

func TestGetUserByEmailToReturnValidUser(t *testing.T) {
	userEmail := "testuser99@example.com"
	userPassword := "password123"
	userID, err := CreateUser(&User{EmailID: userEmail, Password: userPassword})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	user, err := GetUserByEmail(userEmail)
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

func TestGetUserByEmailToReturnErrorForInvalidEmail(t *testing.T) {
	userEmail := "nonexistent@example.com"
	_, err := GetUserByEmail(userEmail)
	if err == nil {
		t.Fatalf("expected GetUserByEmail() to return an error, but got nil")
	}

	expectedErrorMessage := fmt.Sprintf("user with email %s not found", userEmail)
	if err.Error() != expectedErrorMessage {
		t.Errorf("GetUserByEmail() returned wrong error message, got %s, want %s", err.Error(), expectedErrorMessage)
	}
}
