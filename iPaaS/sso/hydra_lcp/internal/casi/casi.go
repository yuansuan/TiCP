package casi

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	//_OAuthBaseURL
	_OAuthBaseURL = "https://auth.ms.casicloud.com/1/oauth/authorize"
	//_OAuthBaseTokenURL
	_OAuthBaseTokenURL = "https://auth.ms.casicloud.com/1/oauth/token"
	//_OAuthBaseGetUserInfoURL
	_OAuthBaseGetUserInfoURL = "https://auth.ms.casicloud.com/1/oauth/user_info"
)

// OAuth CASI oAuth info
type OAuth struct {
	clientID     string
	clientSecret string
	code         string
	token        string
}

// AuthorizeURL get AuthorizeURL
func (o *OAuth) AuthorizeURL(redirectURL string) string {
	u, err := url.Parse(_OAuthBaseURL)
	if err != nil {
		panic(err)
	}

	qs := u.Query()
	qs.Set("redirect_uri", redirectURL)
	qs.Set("client_id", o.clientID)
	qs.Set("response_type", "code")
	qs.Set("ts", strconv.Itoa(int(time.Now().UnixMilli())))
	qs.Set("sign", o.makeSign(qs))
	u.RawQuery = qs.Encode()

	return u.String()
}

// GetTokenURL 获取token方法的url
func (o *OAuth) GetTokenURL(redirectURL string) (string, error) {
	u, err := url.Parse(_OAuthBaseTokenURL)
	if err != nil {
		panic(err)
	}
	//add param
	qs := u.Query()
	qs.Set("grant_type", "authorization_code")
	qs.Set("code", o.code)
	qs.Set("redirect_uri", redirectURL)
	qs.Set("ts", strconv.Itoa(int(time.Now().UnixMilli())))
	qs.Set("client_id", o.clientID)
	qs.Set("sign", o.makeSign(qs))
	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// UserInfo user info
type UserInfo struct {
	Username string `json:"username"`
	Attrs    struct {
		Admin bool `json:"admin"`
	} `json:"attrs"`
}

// GetUserInfoURL get user info url
func (o *OAuth) GetUserInfoURL() (string, error) {
	u, err := url.Parse(_OAuthBaseGetUserInfoURL)
	if err != nil {
	}
	//add param
	qs := u.Query()
	qs.Set("access_token", o.token)
	qs.Set("client_id", o.clientID)
	qs.Set("ts", strconv.Itoa(int(time.Now().UnixMilli())))
	qs.Set("sign", o.makeSign(qs))
	u.RawQuery = qs.Encode()

	return u.String(), nil
}

func (o *OAuth) makeSign(qs url.Values) string {
	var keys []string
	for k := range qs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k + "=" + url.QueryEscape(qs.Get(k)))
	}
	sb.WriteString(o.clientSecret)

	h := md5.New()
	h.Write([]byte(sb.String()))
	return hex.EncodeToString(h.Sum(nil))
}

// NewOAuth build oauth by base
func NewOAuth(clientID, clientSecret string) *OAuth {
	return &OAuth{
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

// NewOAuthCode build oauth by code
func NewOAuthCode(code, clientID, clientSecret string) *OAuth {
	return &OAuth{
		clientID:     clientID,
		clientSecret: clientSecret,
		code:         code,
	}
}

// NewOAuthToken build oauth by token
func NewOAuthToken(token string, clientID string, clientSecret string) *OAuth {
	return &OAuth{
		clientID:     clientID,
		clientSecret: clientSecret,
		token:        token,
	}
}
