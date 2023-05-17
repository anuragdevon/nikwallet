package handlers_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"nikwallet/database"
	"nikwallet/handlers"
	"nikwallet/services"
)

func TestUserHandlers(t *testing.T) {
	db := &database.PostgreSQL{}
	err := db.Connect("testdb")
	if err != nil {
		log.Panic("failed to connect to database:", err)
	}
	defer db.Close()

	userService := services.NewUserService(db.DB)
	authService := services.NewAuthService(db.DB)
	userHandlers := handlers.NewUserHandlers(userService, authService)

	t.Run("SignupHandler to return 201 StatusCreated for valid user creation", func(t *testing.T) {
		signupRequest := map[string]interface{}{
			"email_id": "testhello123@example.com",
			"password": "password123",
		}

		reqBody, err := json.Marshal(signupRequest)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "/signup", bytes.NewReader(reqBody))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(userHandlers.SignupHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)

		var responseBody map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&responseBody)
		assert.NoError(t, err)

		assert.Contains(t, responseBody, "user_id")
		assert.Contains(t, responseBody, "token")
	})

	t.Run("SignupHandler to return 500 InternalServerError for duplicate user entry", func(t *testing.T) {
		// Create a user with the specified email ID using SignupHandler
		newUser := map[string]interface{}{
			"email_id":   "testhello456@example.com",
			"password":   "password123",
			"first_name": "Test",
			"last_name":  "User",
			"profession": "Engineer",
		}

		reqBody, err := json.Marshal(newUser)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "/signup", bytes.NewReader(reqBody))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(userHandlers.SignupHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)

		duplicateUser := map[string]interface{}{
			"email_id":   "testhello456@example.com",
			"password":   "password456",
			"first_name": "Duplicate",
			"last_name":  "User",
			"profession": "Designer",
		}

		reqBody, err = json.Marshal(duplicateUser)
		assert.NoError(t, err)

		req, err = http.NewRequest("POST", "/signup", bytes.NewReader(reqBody))
		assert.NoError(t, err)

		recorder = httptest.NewRecorder()

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

		req, err := http.NewRequest("POST", "/signin", bytes.NewReader(reqBody))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(userHandlers.SigninHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var responseBody map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&responseBody)
		assert.NoError(t, err)

		assert.Contains(t, responseBody, "token")
	})

	t.Run("SigninHandler to return status 401 Unauthorized for invalid user credentials", func(t *testing.T) {
		invalidSigninRequest := map[string]interface{}{
			"email_id": "testhello321@example.com",
			"password": "wrongpassword",
		}

		reqBody, err := json.Marshal(invalidSigninRequest)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "/signin", bytes.NewReader(reqBody))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		http.HandlerFunc(userHandlers.SigninHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

}
