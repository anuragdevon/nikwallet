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

	t.Run("SignupHandler", func(t *testing.T) {
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

	t.Run("SigninHandler", func(t *testing.T) {
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
}
