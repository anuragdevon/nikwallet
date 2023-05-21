package services

import (
	"nikwallet/repository/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	userService := &UserService{
		db: db.DB,
	}
	t.Run("UserService to return error for existing user trying to signup", func(t *testing.T) {
		existingUser := &models.User{
			EmailID:  "existing@example.com",
			Password: "password123",
		}
		_, err := userService.CreateUser(existingUser)

		assert.Nil(t, err)

		duplicateUser := &models.User{
			EmailID:  "existing@example.com",
			Password: "password123",
		}

		noIdCreated, err := userService.CreateUser(duplicateUser)

		assert.Error(t, err, "user already exists")
		assert.Equal(t, 0, noIdCreated)
	})
}
