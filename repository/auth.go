package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var signingKey = []byte("secret-key")

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

func (db *PostgreSQL) AuthenticateUser(email string, password string) (string, error) {
	user, err := db.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if password != user.Password {
		return "", errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: int(user.ID),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	})

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("error generating token: %v", err)
	}

	return tokenString, nil
}

func (db *PostgreSQL) VerifyToken(tokenString string) (*Claims, int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return nil, 0, fmt.Errorf("error parsing token: %v", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, 0, errors.New("error parsing claims")
	}

	if !token.Valid {
		return nil, 0, errors.New("token is invalid")
	}

	return claims, claims.UserID, nil
}
