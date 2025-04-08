package ssojwt

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSsoGenerateParseJwtToken(t *testing.T) {
	secretID := "3XaFonXuzg7"
	secret := []byte("icnEHmh7hIqUP2Lsnu4L")

	assert := require.New(t)
	token, err := GenerateJwtToken(&UserInfo{
		UserID:   "user",
		Fullname: "luke",
	}, secretID, secret, time.Now().Add(10*time.Minute).Unix())

	assert.Nil(err)
	assert.NotEmpty(token)
	fmt.Println(token)

	info, err := ParseJwtToken(token, func(secretID2 string) ([]byte, error) {
		if secretID == secretID2 {
			return secret, nil
		}

		return []byte(""), errors.New("not found")
	})

	assert.Nil(err)
	assert.Equal("user", info.UserID)
	assert.Equal(int64(0), info.CookieExpiredAt)

	_, err = ParseJwtToken(token, func(secretID string) ([]byte, error) {
		return []byte(""), errors.New("not found")
	})

	assert.NotNil(err)

	token, err = GenerateJwtToken(&UserInfo{
		CookieExpiredAt: 10,
	}, secretID, secret, time.Now().Add(10*time.Minute).Unix())

	info, err = ParseJwtToken(token, func(string) ([]byte, error) {
		return secret, nil
	})

	assert.Nil(err)
	assert.Equal(int64(10), info.CookieExpiredAt)
}
