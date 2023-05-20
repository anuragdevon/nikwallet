package services

import (
	"nikwallet/repository/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser_UserExists(t *testing.T) {
	existingUser := &models.User{
		EmailID:  "existing@example.com",
		Password: "password123",
	}

	userService := &UserService{
		db: db.DB,
	}
	createdUserID, err := userService.CreateUser(existingUser)

	assert.Nil(t, err)
	assert.Equal(t, 1, createdUserID)

	duplicateUser := &models.User{
		EmailID:  "existing@example.com",
		Password: "password123",
	}

	noIdCreated, err := userService.CreateUser(duplicateUser)

	assert.Error(t, err, "user already exists")
	assert.Equal(t, 0, noIdCreated)
}
