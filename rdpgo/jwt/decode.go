package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("yuansuan-remote-app")

func Decode(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("parse jwt token failed, %w", err)
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("token invalid: cannot convert to jwt.MapClaims")
	}
	if !token.Valid {
		return "", errors.New("token invalid")
	}

	decodedData, ok := claim["data"].(string)
	if !ok {
		return "", errors.New("token decoded result format not string")
	}

	return decodedData, nil
}
