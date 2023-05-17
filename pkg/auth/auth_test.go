package auth

import (
	"database/sql"
	"log"
	"nikwallet/pkg/db"
	"nikwallet/pkg/user"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestAuthenticateUserWithCorrectCredentials(t *testing.T) {
	email := "testw51@example.com"
	password := "password"

	newUser := &user.User{
		EmailID:  "testw51@example.com",
		Password: "password",
	}

	_, _ = user.CreateUser(newUser)

	token, err := AuthenticateUser(email, password)
	assert.Nil(t, err)
	assert.NotNil(t, token)

	_, _, err = VerifyToken(token)
	assert.Nil(t, err)
}

func TestAuthenticateUserWithIncorrectPassword(t *testing.T) {
	email := "test331@example.com"
	newUser := &user.User{
		EmailID:  "testw51@example.com",
		Password: "test123",
	}

	_, _ = user.CreateUser(newUser)

	token, err := AuthenticateUser(email, "wrong_password")
	assert.NotNil(t, err)
	assert.Equal(t, "", token)
}

func TestAuthenticateUserWithIncorrectEmail(t *testing.T) {
	password := "password"
	newUser := &user.User{
		EmailID:  "testw51@example.com",
		Password: "test123",
	}

	_, _ = user.CreateUser(newUser)

	token, err := AuthenticateUser("wrong_email", password)
	assert.NotNil(t, err)
	assert.Equal(t, "", token)
}
