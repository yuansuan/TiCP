package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func Encode(raw string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": raw,
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("token sign string failed, %w", err)
	}

	return tokenString, nil
}
