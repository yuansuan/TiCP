package handler

import (
	"context"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/common"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/with"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"
)

// SignupAndLogin SignupAndLogin
// swagger:route POST /api/signupandlogin fe SignupReq
// @POST /api/signupandlogin
func (h *Handler) SignupAndLogin(c *gin.Context) {
	//1.注册
	var req SignupAndLoginReq
	if err := c.BindJSON(&req); err != nil {
		http.Errf(c, consts.ErrHydraLcpBadRequest, "err: %v", err)
		return
	}

	// keyid 存在，解密
	password := req.Password
	if req.KeyID != "" && req.Password != "" {
		cache := boot.MW.DefaultCache()
		var loginKey string
		_, ok := cache.Get(common.HydraLcpLoginKey, req.KeyID, &loginKey)
		if !ok {
			http.Err(c, consts.ErrHydraLcpLoginKeyExpire, "login key is expired")
			return
		}
		cipherPwd, err := base64.StdEncoding.DecodeString(req.Password)
		if err != nil {
			http.Err(c, consts.ErrHydraLcpLoginKeyExpire, "invalid base64 encoding")
			return
		}
		pwd, err := common.AESDecrypt(cipherPwd, []byte(loginKey))

		if err != nil {
			http.Err(c, consts.ErrHydraLcpLoginKeyExpire, "failed to decrypt")
			return
		}
		password = pwd
	}

	re := util.IsPasswordValid(password)
	if !re {
		http.Errf(c, consts.ErrHydraLcpPwdInvalidate, "password is invalid")
		return
	}

	//校验 验证码
	if req.SignupType == common.TypePhone {
		err := h.phoneSrv.VerifyCode(c, req.Phone, req.Code)
		if err != nil {
			http.ErrFromGrpc(c, err)
			return
		}
	} else {
		err := h.emailSrv.VerifyCode(c, req.Email, req.Code)
		if err != nil {
			http.ErrFromGrpc(c, err)
			return
		}
	}

	if req.UserSource == "" {
		req.UserSource = "未知来源"
	}
	if req.UserChannel == "" {
		req.UserChannel = "UNKNOWN"
	}

	var redirectURL, userId string
	var code int
	err := with.DefaultTransaction(c, func(context context.Context) (err error) {
		//1.get userID, if user not exist, add User
		_, err = h.addUser(context, req.Phone, "", password, req.Name, req.Company, req.Email, req.UserSource, req.UserChannel)
		if err != nil {
			code = int(status.Code(err))
			return err
		}
		//2.登录
		redirectURL, userId, code, err = h.AutomaticLogin(context, req)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		http.Err(c, codes.Code(code), err.Error())
		return
	}

	if userId != "" {
		http.Ok(c, userId)
		return
	}
	http.Ok(c, LoginResp{RedirectURL: redirectURL})
}

// SignupAndLoginReqResp SignupAndLoginReqResp
// swagger:response SignupAndLoginReqResp
type SignupAndLoginReqResp struct {
	UserID string `json:"user_id"`
}

// SignupAndLoginReq SignupAndLoginReq
// swagger:parameters SignupAndLoginReq
type SignupAndLoginReq struct {
	Phone               string           `json:"phone"`
	Code                string           `json:"code"`
	Password            string           `json:"password"`
	UserSource          string           `json:"user_source"`
	UserChannel         string           `json:"user_channel"`
	Company             string           `json:"company"`
	Email               string           `json:"email"`
	Name                string           `json:"name"`
	WithHydra           bool             `json:"with_hydra"`
	Challenge           string           `json:"challenge"`
	LoginType           common.LoginType `json:"login_type"`
	KeyID               string           `json:"key_id"`
	LdapCN              string           `json:"ldap_cn"`
	ImageCaptchaID      string           `json:"image_captcha_id"`
	ImageCaptchaContent string           `json:"image_captcha_content"`
	SignupType          string           `json:"signup_type"`
}
