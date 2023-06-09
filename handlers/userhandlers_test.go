package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"nikwallet/repository/models"
	"nikwallet/services"
)

func TestUserHandlers(t *testing.T) {

	userService := services.NewUserService(db.DB)
	authService := services.NewAuthService(db.DB)
	userHandlers := NewUserHandlers(userService, authService)

	t.Run("SignupHandler to return 201 StatusCreated for valid user creation", func(t *testing.T) {
		signupRequest := map[string]interface{}{
			"email_id": "testhello123@example.com",
			"password": "password123",
		}

		reqBody, err := json.Marshal(signupRequest)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "user/signup", bytes.NewReader(reqBody))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(userHandlers.SignupHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)

		var responseBody map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&responseBody)
		assert.NoError(t, err)

		assert.Contains(t, responseBody, "user_id")
		assert.Contains(t, responseBody, "id_token")
	})

	t.Run("SignupHandler to return 500 InternalServerError for duplicate user entry", func(t *testing.T) {
		newUser := &models.User{
			EmailID:  "nikwallethello@example.com",
			Password: "password",
		}

		_, err := userService.CreateUser(newUser)
		assert.NoError(t, err)

		duplicateUser := map[string]interface{}{
			"email_id": "nikwallethello@example.com",
			"password": "password4561",
		}

		reqBody, err := json.Marshal(duplicateUser)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "user/signup", bytes.NewReader(reqBody))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(userHandlers.SignupHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

	t.Run("SignInHandler to return status 200 StatusOk for successful user login", func(t *testing.T) {
		signinRequest := map[string]interface{}{
			"email_id": "testhello321@example.com",
			"password": "password123",
		}

		reqBody, err := json.Marshal(signinRequest)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "user/signin", bytes.NewReader(reqBody))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(userHandlers.SigninHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var responseBody map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&responseBody)
		assert.NoError(t, err)

		assert.Contains(t, responseBody, "id_token")
	})

	t.Run("SigninHandler to return status 401 Unauthorized for invalid user credentials", func(t *testing.T) {
		invalidSigninRequest := map[string]interface{}{
			"email_id": "emaildoesnotexits@example.com",
			"password": "wrongpassword",
		}

		reqBody, err := json.Marshal(invalidSigninRequest)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "user/signin", bytes.NewReader(reqBody))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(userHandlers.SigninHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})
}
