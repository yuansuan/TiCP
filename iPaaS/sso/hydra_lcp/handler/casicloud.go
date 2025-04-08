package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ory/hydra/sdk/go/hydra/client/admin"
	"github.com/ory/hydra/sdk/go/hydra/models"
	"google.golang.org/grpc/codes"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/config"
	localModels "github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/internal/casi"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/rpc"

	yuansuanhttp "github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
)

// CASIUserRes casi user info
type CASIUserRes struct {
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
		ExpiresIn    string `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
		UserOpenID   string `json:"user_open_id"`
		ClientID     string `json:"client_id"`
	} `json:"data"`
}

// CallCASI request with CASI
func (handler *Handler) CallCASI(c *gin.Context) {
	clientID := config.Custom.CasiOauth.ClientID
	clientSecret := config.Custom.CasiOauth.ClientSecret
	redirectURL := config.Custom.CasiOauth.RedirectURL
	loginChallenge, ok := c.GetQuery("login_challenge")
	if !ok {
		c.JSON(500, nil)
		return
	}
	oauth := casi.NewOAuth(clientID, clientSecret)
	redirectURL = redirectURL + "?" + "login_challenge=" + loginChallenge
	c.Redirect(http.StatusFound, oauth.AuthorizeURL(redirectURL))
}

// CASICallback callback
func (handler *Handler) CASICallback(c *gin.Context) {
	redirectURL := config.Custom.CasiOauth.RedirectURL
	clientID := config.Custom.CasiOauth.ClientID
	clientSecret := config.Custom.CasiOauth.ClientSecret
	//1、get code and login_challenge
	code, ok := c.GetQuery("code")
	if !ok {
		yuansuanhttp.Errf(c, consts.ErrHydraLcpBadRequest, "get code error")
		return
	}
	loginChallenge, ok := c.GetQuery("login_challenge")
	if !ok {
		yuansuanhttp.Errf(c, consts.ErrHydraLcpBadRequest, "get  login_challenge error")
		return
	}
	//2、build request
	request := casi.NewRequestBase(code, clientID, clientSecret, redirectURL)
	//3、get casi token by use code
	authToken, err := request.GetCASIOauthToken()
	if err != nil {
		yuansuanhttp.Errf(c, consts.ErrHydraLcpBadRequest, "get casi token error: %s", err)
		return
	}
	//4、get user info by use token
	request.Token = authToken.Data.AccessToken
	casiUserInfo, err := request.GetCASIOauthUserInfo()
	if err != nil {
		yuansuanhttp.Errf(c, consts.ErrHydraLcpBadRequest, "get casi user info error: %s", err)
		return
	}
	//5、init ssoUser
	userModel := &localModels.SsoUser{Phone: casiUserInfo.Data.UserMobile}
	exist, err := handler.userSrv.Get(c, userModel)
	if err != nil {
		yuansuanhttp.Errf(c, codes.Internal, "get sso user error: %s", err)
		return
	}
	var userID int64
	if !exist {
		newID, err := rpc.GenID(c)
		if err != nil {
			yuansuanhttp.Errf(c, consts.ErrHydraLcpBadRequest, "get rpc id error: %s", err)
			return
		}
		userID = int64(newID)
		err = handler.userSrv.Add(c, &localModels.SsoUser{
			Ysid:        userID,
			Name:        "航天云网用户",
			RealName:    "航天云网用户",
			Phone:       casiUserInfo.Data.UserMobile,
			UserChannel: "航天云网",
			UserSource:  "航天云网",
			UserReferer: 1,
			CreateTime:  time.Now(),
			ModifyTime:  time.Now(),
			IsActivated: true,
		}, "")
		if err != nil {
			yuansuanhttp.Errf(c, consts.ErrHydraLcpBadRequest, "add sso-user info: %s", err)
			return
		}
	} else {
		userID = userModel.Ysid
	}
	sid := snowflake.ID(userID).String()
	//6、login yuansuan
	res, err := handler.HydraClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
		WithLoginChallenge(loginChallenge).
		WithBody(&models.HandledLoginRequest{
			Subject:     &sid,
			Remember:    true,
			RememberFor: 1000,
		}))
	if err != nil {
		yuansuanhttp.Errf(c, consts.ErrHydraLcpFailedToReqHydra, "unable to send request for login to hydra: %v", err)
		return
	}
	c.Redirect(http.StatusFound, res.Payload.RedirectTo)
}
