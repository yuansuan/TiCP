package v1

import (
	"encoding/base64"
	"fmt"
	jwtgo "github.com/golang-jwt/jwt/v4"
	"math/rand"
	"strings"
)

const (
	// Minimum length for MinIO access key.
	accessKeyMinLen = 3

	// Maximum length for MinIO access key.
	// There is no max length enforcement for access keys
	accessKeyMaxLen = 20

	// Minimum length for MinIO secret key for both server
	secretKeyMinLen = 8

	// Maximum secret key length for MinIO, this
	// is used when autogenerating new credentials.
	// There is no max length enforcement for secret keys
	secretKeyMaxLen = 40

	// Alpha numeric table used for generating access keys.
	alphaNumericTable = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Total length of the alpha numeric table.
	alphaNumericTableLen = byte(len(alphaNumericTable))

	// key for JWT token
	// FIXME: read from config file
	globalSecretKey = "1234yskj"
)

func GenerateCredentials() (accessKey, secretKey string, err error) {
	readBytes := func(size int) (data []byte, err error) {
		data = make([]byte, size)
		var n int
		if n, err = rand.Read(data); err != nil {
			return nil, err
		} else if n != size {
			return nil, fmt.Errorf("Not enough data. Expected to read: %v bytes, got: %v bytes", size, n)
		}
		return data, nil
	}

	// Generate access key.
	keyBytes, err := readBytes(accessKeyMaxLen)
	if err != nil {
		return "", "", err
	}
	for i := 0; i < accessKeyMaxLen; i++ {
		keyBytes[i] = alphaNumericTable[keyBytes[i]%alphaNumericTableLen]
	}
	accessKey = string(keyBytes)

	// Generate secret key.
	keyBytes, err = readBytes(secretKeyMaxLen)
	if err != nil {
		return "", "", err
	}

	secretKey = strings.ReplaceAll(string([]byte(base64.StdEncoding.EncodeToString(keyBytes))[:secretKeyMaxLen]),
		"/", "+")

	return accessKey, secretKey, nil
}

func GetNewCredentialsWithMetadata(m map[string]interface{}) (accessKey, secretKey, sessionToken string, err error) {
	accessKey, secretKey, err = GenerateCredentials()
	if err != nil {
		return "", "", "", err
	}
	m["accessKey"] = accessKey
	sessionToken, err = JWTSignWithAccessKey(accessKey, m, globalSecretKey)
	if err != nil {
		return "", "", "", err
	}
	return accessKey, secretKey, sessionToken, nil
}

func JWTSignWithAccessKey(accessKey string, m map[string]interface{}, tokenSecret string) (string, error) {
	m["accessKey"] = accessKey
	jwt := jwtgo.NewWithClaims(jwtgo.SigningMethodHS512, jwtgo.MapClaims(m))
	return jwt.SignedString([]byte(tokenSecret))
}
