package casi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// UserRes casi user info
type UserRes struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		IsAdmin     bool   `json:"is_admin"`
		UserMobile  string `json:"user_mobile"`
		UserName    string `json:"user_name"`
		UserOpenID  string `json:"user_open_id"`
		UserAccount string `json:"user_account"`
		OrgName     string `json:"org_name"`
		OrgOpenID   string `json:"org_open_id"`
	} `json:"data"`
}

// AuthTokenRes casi token info
type AuthTokenRes struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		AccessToken string `json:"access_token"`
		//access token expire time 8 hour,unit second
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
		UserOpenID   string `json:"user_open_id"`
		ClientID     string `json:"client_id"`
	} `json:"data"`
}

// RequestBase request info
type RequestBase struct {
	Code         string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Token        string
}

// NewRequestBase build request base
func NewRequestBase(code, clientID, clientSecret, redirectURL string) *RequestBase {
	return &RequestBase{
		Code:         code,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
	}
}

// GetCASIOauthToken get token by code
func (request *RequestBase) GetCASIOauthToken() (*AuthTokenRes, error) {
	oauth := NewOAuthCode(request.Code, request.ClientID, request.ClientSecret)
	url, err := oauth.GetTokenURL(request.RedirectURL)
	if err != nil {
		return &AuthTokenRes{}, err
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return &AuthTokenRes{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &AuthTokenRes{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte
	if err != nil {
		return &AuthTokenRes{}, err
	}
	var authCode AuthTokenRes
	err = json.Unmarshal(body, &authCode)
	if err != nil {
		return &AuthTokenRes{}, err
	}
	if authCode.Code != 200 {
		return &AuthTokenRes{}, err
	}
	return &authCode, nil
}

// GetCASIOauthUserInfo get user info by token
func (request *RequestBase) GetCASIOauthUserInfo() (*UserRes, error) {
	oauth := NewOAuthToken(request.Token, request.ClientID, request.ClientSecret)
	url, err := oauth.GetUserInfoURL()
	if err != nil {
		return &UserRes{}, err
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return &UserRes{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &UserRes{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &UserRes{}, err
	}
	var userInfo UserRes
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return &UserRes{}, err
	}

	if userInfo.Code != 200 {
		return &UserRes{}, err
	}
	return &userInfo, nil
}
