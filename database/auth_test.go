package database

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	db := &PostgreSQL{}
	err := db.Connect("testdb")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	t.Run("Authenticate method to authenticate user with correct credentials", func(t *testing.T) {
		email := "testw51@example.com"
		password := "password"

		newUser := &User{
			EmailID:  "testw51@example.com",
			Password: "password",
		}

		_, _ = db.CreateUser(newUser)

		token, err := db.AuthenticateUser(email, password)
		assert.Nil(t, err)
		assert.NotNil(t, token)

		_, _, err = db.VerifyToken(token)
		assert.Nil(t, err)
	})

	t.Run("Authenticate method to return error for invalid password", func(t *testing.T) {
		email := "test331@example.com"
		newUser := &User{
			EmailID:  "testw51@example.com",
			Password: "test123",
		}

		_, _ = db.CreateUser(newUser)

		token, err := db.AuthenticateUser(email, "wrong_password")
		assert.NotNil(t, err)
		assert.Equal(t, "", token)
	})
	t.Run("Authenticate method to return error for invalid email", func(t *testing.T) {
		password := "password"
		newUser := &User{
			EmailID:  "testw51@example.com",
			Password: "test123",
		}

		_, _ = db.CreateUser(newUser)

		token, err := db.AuthenticateUser("wrong_email", password)
		assert.NotNil(t, err)
		assert.Equal(t, "", token)
	})

	t.Run("VerifyToken method to successfully verify valid token", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
			UserID: 123,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			},
		})

		tokenString, _ := token.SignedString(signingKey)

		claims, userID, err := db.VerifyToken(tokenString)

		assert.Nil(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, 123, userID)
	})

	t.Run("VerifyToken method to return error for invalid token key", func(t *testing.T) {
		invalidKey := []byte("invalid key")
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
			UserID: 123,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			},
		})

		tokenString, _ := token.SignedString(invalidKey)

		claims, userID, err := db.VerifyToken(tokenString)

		assert.NotNil(t, err)
		assert.Nil(t, claims)
		assert.Equal(t, 0, userID)
	})

	t.Run("VerifyToken method to return error for expired token key", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
			UserID: 123,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(-24 * time.Hour).Unix(),
			},
		})

		tokenString, _ := token.SignedString(signingKey)

		claims, userID, err := db.VerifyToken(tokenString)

		assert.NotNil(t, err)
		assert.Nil(t, claims)
		assert.Equal(t, 0, userID)
	})
}
