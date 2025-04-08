package util

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"

	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/common"
)

// JWTGenerate JWTGenerate
func JWTGenerate(prefix string, subject string, expireAt time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   prefix + "-" + subject,
		IssuedAt:  time.Now().UTC().Unix(),
		NotBefore: time.Date(2019, 9, 29, 0, 0, 0, 0, time.UTC).Unix(),
		ExpiresAt: expireAt.UTC().Unix(),
	})

	// Sign and get the complete encoded tokenString as a string using the secret
	tokenString, err := token.SignedString(common.JWTSecret)
	if err != nil {
		return "", status.Errorf(consts.ErrHydraLcpJWTGenerate, "fail to generate jwt for %v, err: %v", prefix+subject, err)
	}

	return tokenString, nil
}

// JWTParse JWTParse
func JWTParse(tokenString string) (jwt.StandardClaims, error) {
	var claims jwt.StandardClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return common.JWTSecret, nil
	})

	if err != nil {
		return claims, status.Errorf(consts.ErrHydraLcpJWTInvalid, "token invalid, err: %v", err)
	}

	if !token.Valid {
		return claims, status.Error(consts.ErrHydraLcpJWTInvalid, "token invalid")
	}

	return claims, nil
}

// JWTGetSubject JWTGetSubject
func JWTGetSubject(prefix string, tokenString string) (string, error) {
	claims, err := JWTParse(tokenString)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(claims.Subject, prefix+"-"), nil
}
