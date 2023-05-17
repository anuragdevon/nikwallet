package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	user "nikwallet/pkg/user"
)

var signingKey = []byte("secret-key")

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

func AuthenticateUser(email string, password string) (string, error) {
	u, err := user.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if password != u.Password {
		return "", errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: u.ID,
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

func VerifyToken(tokenString string) (*Claims, int, error) {
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
