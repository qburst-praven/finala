package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte("your-super-secret-and-long-enough-key-please-change-this-ASAP")

// GenerateJWT creates a new JWT for a given username.

func GenerateJWT(username string) (string, time.Time, error) {

	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours

	claims := &jwt.RegisteredClaims{

		Subject: username,

		ExpiresAt: jwt.NewNumericDate(expirationTime),

		IssuedAt: jwt.NewNumericDate(time.Now()),

		Issuer: "finala-api",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecretKey)

	if err != nil {

		return "", time.Time{}, err

	}

	return tokenString, expirationTime, nil

}
