package auth

import (
	"database/sql"
	"log"
	"nikwallet/pkg/db"
	"nikwallet/pkg/user"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
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

func TestVerifyTokenWithValidToken(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: 123,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	})

	tokenString, _ := token.SignedString(signingKey)

	claims, userID, err := VerifyToken(tokenString)

	assert.Nil(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, 123, userID)
}

func TestVerifyTokenWithInvalidToken(t *testing.T) {
	invalidKey := []byte("invalid key")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: 123,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	})

	tokenString, _ := token.SignedString(invalidKey)

	claims, userID, err := VerifyToken(tokenString)

	assert.NotNil(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, 0, userID)
}

func TestVerifyTokenWithExpiredToken(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: 123,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-24 * time.Hour).Unix(),
		},
	})

	tokenString, _ := token.SignedString(signingKey)

	claims, userID, err := VerifyToken(tokenString)

	assert.NotNil(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, 0, userID)
}
