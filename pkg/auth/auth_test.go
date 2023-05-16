package auth

import (
	"nikwallet/pkg/db"
	user "nikwallet/pkg/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticateUserWithCorrectCredentials(t *testing.T) {
	email := "testw51@example.com"
	password := "password"
	db, err := db.ConnectToDB("testdb")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	newUser := &user.User{
		EmailID:  "testw51@example.com",
		Password: "password",
	}

	_, _ = user.CreateUser(db, newUser)

	token, err := AuthenticateUser(db, email, password)
	assert.Nil(t, err)
	assert.NotNil(t, token)

	_, _, err = VerifyToken(token)
	assert.Nil(t, err)
}

func TestAuthenticateUserWithIncorrectPassword(t *testing.T) {
	email := "test331@example.com"
	db, err := db.ConnectToDB("testdb")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	newUser := &user.User{
		EmailID:  "testw51@example.com",
		Password: "test123",
	}

	_, _ = user.CreateUser(db, newUser)

	token, err := AuthenticateUser(db, email, "wrong_password")
	assert.NotNil(t, err)
	assert.Equal(t, "", token)
}

func TestAuthenticateUserWithIncorrectEmail(t *testing.T) {
	password := "password"
	db, err := db.ConnectToDB("testdb")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	newUser := &user.User{
		EmailID:  "testw51@example.com",
		Password: "test123",
	}

	_, _ = user.CreateUser(db, newUser)

	token, err := AuthenticateUser(db, "wrong_email", password)
	assert.NotNil(t, err)
	assert.Equal(t, "", token)
}
