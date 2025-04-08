package ssojwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// UserInfo holds the info of user
type UserInfo struct {
	// CookieExpiredAt indicates the expiration time of cookie (int unix seconds)
	// 0 means cookies expired at the end of session (when user closed the browser)
	CookieExpiredAt int64 `json:"cookie_expired_at"`

	// UserID field indicates the user id of login user
	UserID string `json:"user_id"`

	// Fullname field indicates the fullname of login user, optional
	Fullname string `json:"fullname"`

	// Email optional
	Email string `json:"email"`

	// Phone optional
	Phone string `json:"phone"`
}

// JwtClaims holds the k-v to issue a cookie
type JwtClaims struct {
	*UserInfo

	// The Issuer field holds the id of secret
	// Center should use the id of secret to get the secret which is used to
	// verify the jwt token
	*jwt.StandardClaims
}

// GenerateJwtToken generates the jwt token to exchange cookie, expireAt need unix time, eg time.Now().Add(10 * time.Minute).Unix())
func GenerateJwtToken(userInfo *UserInfo, secretID string, secret []byte, expiresAt int64) (token string, err error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{
		StandardClaims: &jwt.StandardClaims{
			Audience:  secretID,
			ExpiresAt: expiresAt,
			IssuedAt:  time.Now().Unix(),
		},

		UserInfo: userInfo,
	})

	return t.SignedString(secret)
}

// GetSsoJwtSecret the function to retrieve secret with secretID
type GetSsoJwtSecret func(secretID string) ([]byte, error)

func getStringAttr(attrs map[string]interface{}, attr string) string {
	if v, ok := attrs[attr]; ok {
		if value, ok := v.(string); ok {
			return value
		}
	}

	return ""
}

// ParseJwtToken parses the jwt token
func ParseJwtToken(tokenString string, getSecret GetSsoJwtSecret) (u *UserInfo, err error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		claims := t.Claims.(jwt.MapClaims)
		return getSecret(claims["aud"].(string))
	})

	if err != nil {
		return nil, err
	}

	return &UserInfo{
		CookieExpiredAt: int64(token.Claims.(jwt.MapClaims)["cookie_expired_at"].(float64)),
		UserID:          getStringAttr(token.Claims.(jwt.MapClaims), "user_id"),
		Fullname:        getStringAttr(token.Claims.(jwt.MapClaims), "fullname"),
		Email:           getStringAttr(token.Claims.(jwt.MapClaims), "email"),
		Phone:           getStringAttr(token.Claims.(jwt.MapClaims), "phone"),
	}, nil
}
