package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"io/ioutil"
	http2 "net/http"
	"net/url"
	"os"
	"sync"
	"time"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	"github.com/gin-gonic/gin"

	"github.com/ory/hydra/sdk/go/hydra/client/admin"
	"github.com/ory/hydra/sdk/go/hydra/models"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/iPaaS/sso/protos/platform/idgen"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/common"
	localModels "github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/service"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

var mtx sync.Mutex

// HydraLogin HydraLogin
// swagger:route GET /api/hydra/login  hydra
//
// login api used by hydra, it's behavior is strictly defined by hydra.
//
//	    Responses:
//		     302
//	      90001: ErrHydraLcpFailedToReqHydra
//
// @GET /api/hydra/login
func (h *Handler) HydraLogin(c *gin.Context) {
	challenge := c.Query(common.HydraLoginChallenge)
	logger := logging.GetLogger(c)
	logger.Infof("[hydra login] start hydra login for challenge %v", challenge)
	logger.Info(">>>>>>>>>>>>>>> challenge: ", challenge)

	// get login request from hydra
	res, err := h.HydraClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(challenge))
	if err != nil {
		http.Errf(c, consts.ErrHydraLcpFailedToReqHydra, "unable to send request for login to hydra: %v", err)
		return
	}

	logger.Info(">>>>>>>>>>>>>>> response: ", res)

	//If hydra was already able to authenticate the user, skip will be true and we do not need to re-authenticate the user.
	if res.Payload.Skip {
		res, err := h.HydraClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
			WithLoginChallenge(challenge).
			WithBody(&models.HandledLoginRequest{
				Subject:     &res.Payload.Subject,
				Remember:    true,
				RememberFor: h.TokenExpireTime,
			}))
		if err != nil {
			logger.Warnf("[hydra login exception] unable to send request for login to hydra: %v", err)
			http.Errf(c, consts.ErrHydraLcpFailedToReqHydra, "unable to send request for login to hydra: %v", err)
			return
		}
		c.Redirect(http2.StatusFound, res.Payload.RedirectTo)
	}

	logger.Infof("[hydra login] redirect to %v", h.getLoginURL(c, challenge))
	c.Redirect(http2.StatusFound, h.getLoginURL(c, challenge))
	logger.Infof("[hydra login] hydra login successful for challenge %v", challenge)
}

// HydraPortalLogin HydraPortalLogin
// swagger:route GET /api/hydra_portal/login  hydra_portal HydraLoginReq
//
// login api used by hydra-portal, it's behavior is strictly defined by hydra-portal.
//
//	    Responses:
//		     302
//	      90001: ErrHydraLcpFailedToReqHydra
//
// @GET /api/hydra_portal/login
func (h *Handler) HydraPortalLogin(c *gin.Context) {
	challenge := c.Query(common.HydraLoginChallenge)
	logger := logging.GetLogger(c)
	logger.Infof("[hydra-portal login] start hydra-portal login for challenge %v", challenge)
	logger.Info(">>>>>>>>>>>>>>> challenge: ", challenge)

	// get login request from hydra-portal
	res, err := h.HydraClientPortal.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(challenge))
	if err != nil {
		logger.Warnf("[hydra-portal login exception] unable to send request for login to hydra-portal: %v", err)
		http.Errf(c, consts.ErrHydraLcpFailedToReqHydra, "unable to send request for login to hydra-portal: %v", err)
		return
	}

	logger.Info(">>>>>>>>>>>>>>> response: ", res)

	//If hydra-portal was already able to authenticate the user, skip will be true and we do not need to re-authenticate the user.
	if res.Payload.Skip {
		res, err := h.HydraClientPortal.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
			WithLoginChallenge(challenge).
			WithBody(&models.HandledLoginRequest{
				Subject:     &res.Payload.Subject,
				Remember:    true,
				RememberFor: h.TokenExpireTimePortal,
			}))
		if err != nil {
			logger.Warnf("[hydra-portal login exception] unable to send request for login to hydra-portal: %v", err)
			http.Errf(c, consts.ErrHydraLcpFailedToReqHydra, "unable to send request for login to hydra-portal: %v", err)
			return
		}
		c.Redirect(http2.StatusFound, res.Payload.RedirectTo)
	}

	logger.Infof("[hydra-portal login] redirect to %v", h.getLoginURLPortal(c, challenge))
	c.Redirect(http2.StatusFound, h.getLoginURLPortal(c, challenge))
	logger.Infof("[hydra-portal login] hydra-portal login successful for challenge %v", challenge)
}

// HydraLoginReq HydraLoginReq
// swagger:parameters HydraLoginReq
type HydraLoginReq struct {
	// in: query
	LoginChallenge string `json:"login_challenge"`
}

var mutexSendVerificationCode sync.Mutex

// SendVerificationCode SendVerificationCode
// swagger:route GET /api/send_code fe SendVerificationCodeReq
//
// Responses:
//
//	200
//	90002: ErrHydraLcpBadRequest
//	90027：ErrHydraLcpPhoneExist
//	90035: ErrHydraLcpPhoneNotExist
//
// @GET /api/send_code
// @example:
//
//	curl http://127.0.0.1:8899/api/send_code?phone=15900000000&forget=1
//	curl http://127.0.0.1:8899/api/send_code?phone=15900000000
//	curl curl http://127.0.0.1:8899/api/send_code?phone=15900000000&login=1
func (h *Handler) SendVerificationCode(c *gin.Context) {
	phone := c.Query("phone")
	if phone == "" {
		http.Err(c, consts.ErrHydraLcpBadRequest, "lack phone number")
		return
	}

	mutexSendVerificationCode.Lock()
	defer mutexSendVerificationCode.Unlock()

	cache := boot.MW.DefaultCache()
	logger := logging.GetLogger(c)

	// check whether send count over the max(per hour&IP ).
	ipAdress := c.ClientIP()
	var sendCount int64
	_, ok := cache.Get(common.RedisPrefixSendCountPerHourAndIP, ipAdress, &sendCount)
	if ok {

		if sendCount >= h.PerIPHourSendPhoneCodeMax {
			logger.Warnf("[send code exception] send count over; phone : %v, IP: %v ", phone, ipAdress)
			http.Err(c, consts.ErrHydraLcpWait, "send code fail to phone "+phone)
			return
		}
	}

	sendCount++
	cache.PutWithExpire(common.RedisPrefixSendCountPerHourAndIP, ipAdress, sendCount, time.Minute*60)
	//跳过手机号检查 用于营销活动一键登录/注册

	switch {
	case c.Query("forget") == "1" || c.Query("login") == "1": //
		exists, err := h.userSrv.Get(c, &localModels.SsoUser{Phone: phone})
		if err != nil {
			http.ErrFromGrpc(c, err)
			return
		}
		if !exists {
			http.Err(c, consts.ErrHydraLcpPhoneNotExist, "phone number not exist")
			return
		}
	case c.Query("skip_check") == "1":
		// do nothing
	default: // register
		exists, err := h.userSrv.Get(c, &localModels.SsoUser{Phone: phone})
		if err != nil {
			http.ErrFromGrpc(c, err)
			return
		}
		if exists {
			http.Err(c, consts.ErrHydraLcpPhoneExist, "phone number is existed")
			return
		}
	}

	// 切换签名
	sign := service.SignYS
	if c.Query("pid") == consts.ZSWLProductID.String() {
		sign = service.SignZS
	}
	err := h.phoneSrv.SendVerificationCode(c, phone, sign)
	if err != nil {
		http.ErrFromGrpc(c, err)
		return
	}

	http.Ok(c, nil)
}

func (h *Handler) NewSendVerificationCode(c *gin.Context) {
	recipient := c.Query("recipient")
	if recipient == "" {
		http.Err(c, consts.ErrHydraLcpBadRequest, "lack recipient field")
		return
	}
	sendType := c.Query("sendType")
	if sendType == "" {
		http.Err(c, consts.ErrHydraLcpBadRequest, "lack sendType field")
		return
	}

	mutexSendVerificationCode.Lock()
	defer mutexSendVerificationCode.Unlock()

	cache := boot.MW.DefaultCache()
	logger := logging.GetLogger(c)

	// check whether send count over the max(per hour&IP ).
	ipAdress := c.ClientIP()
	var sendCount int64
	_, ok := cache.Get(common.RedisPrefixSendCountPerHourAndIP, ipAdress, &sendCount)
	if ok {

		if sendCount >= h.PerIPHourSendPhoneCodeMax {
			logger.Warnf("[send code exception] send count over; %v, IP: %v ", recipient, ipAdress)
			http.Err(c, consts.ErrHydraLcpWait, "send code fail to "+recipient)
			return
		}
	}

	sendCount++
	cache.PutWithExpire(common.RedisPrefixSendCountPerHourAndIP, ipAdress, sendCount, time.Minute*60)
	//跳过手机号检查 用于营销活动一键登录/注册
	var code int
	var message string
	switch {
	case c.Query("forget") == "1" || c.Query("login") == "1":

		user := localModels.SsoUser{}
		if sendType == "email" {
			user.Email = recipient
			code = consts.ErrHydraLcpEmailNotExist
			message = "email not exist"
		}
		if sendType == "phone" {
			user.Phone = recipient
			code = consts.ErrHydraLcpPhoneNotExist
			message = "phone number not exist"
		}

		exists, err := h.userSrv.Get(c, &user)
		if err != nil {
			http.ErrFromGrpc(c, err)
			return
		}
		if !exists {
			http.Err(c, codes.Code(code), message)
			return
		}
	case c.Query("skip_check") == "1":
		// do nothing
	default: // register
		user := localModels.SsoUser{}
		if sendType == common.TypeEmail {
			user.Phone = recipient
			code = consts.ErrHydraLcpEmailExist
			message = "email is existed"
		}
		if sendType == common.TypePhone {
			user.Email = recipient
			code = consts.ErrHydraLcpPhoneExist
			message = "phone number is existed"
		}
		exists, err := h.userSrv.Get(c, &user)
		if err != nil {
			http.ErrFromGrpc(c, err)
			return
		}
		if exists {
			http.Err(c, codes.Code(code), message)
			return
		}
	}

	// 切换签名
	sign := service.SignYS
	if c.Query("pid") == consts.ZSWLProductID.String() {
		sign = service.SignZS
	}

	// 发送验证码（手机或者邮箱）
	if sendType == common.TypeEmail {
		err := h.emailSrv.SendVerificationCode(c, recipient)
		if err != nil {
			http.ErrFromGrpc(c, err)
			return
		}
	} else {
		err := h.phoneSrv.SendVerificationCode(c, recipient, sign)
		if err != nil {
			http.ErrFromGrpc(c, err)
			return
		}
	}

	http.Ok(c, nil)
}

// SendVerificationCodeReq SendVerificationCodeReq
// swagger:parameters SendVerificationCodeReq
type SendVerificationCodeReq struct {
	// required: true
	// in: query
	Phone string `json:"phone"`
}

// ForwardToDefault ForwardToDefault
// swagger:route GET /api/forward_to_default fe nil
//
// this api is currently used for to forward to default destination when error occurs
//
// forward to default destination when error occurs
//
//	    Responses:
//				200: ForwardToDefaultResp
//
// @GET /api/forward_to_default
func (h *Handler) ForwardToDefault(c *gin.Context) {
	logger := logging.GetLogger(c)
	defaultPage := os.Getenv("DEFAULT_PAGE")
	if defaultPage == "" {
		defaultPage = "https://www.yuansuan.cn/login"
		logger.Warnf("[api login exception] env for DEFAULT_PAGE not found, use %v as default page ", defaultPage)
	} else {
		logger.Info("[api login] env DEFAULT_PAGE is ", defaultPage)
	}
	http.Ok(c, ForwardToDefaultResp{DefaultURL: defaultPage})
}

// Login Login
// swagger:route POST /api/login fe LoginReq
//
// this api is currently used for login with phone+code and phone+password
//
// login with phone and password
//
//			login type: 2
//			req: phone, password
//
//	    Responses:
//				200: LoginResp
//				90001: ErrHydraLcpFailedToReqHydra
//				90002: ErrHydraLcpBadRequest
//	         90029: ErrHydraLcpCaptchaVerifyFailed
//	         90038: ErrHydraLcpUserLoginFreezed
//	         90040: ErrHydraLcpLoginFailed
//	         90041: ErrHydraLcpLoginKeyExpire
//	         90042: ErrHydraLcpGetLoginKeyFail
//
// @POST /api/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.BindJSON(&req); err != nil {
		http.Errf(c, consts.ErrHydraLcpBadRequest, "err: %v", err)
		return
	}

	var userID int64
	var err error
	logger := logging.GetLogger(c)

	// keyid 存在，解密
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
		req.Password = string(pwd)
	}

	mtx.Lock()
	defer mtx.Unlock()

	// check username and password
	switch req.LoginType {
	case common.EmailPassword:
		// 判断是否超过错误次数
		isOverMax := h.checkVerifyFailCountOverMax(c, req.Email)
		if isOverMax {
			logger.Infof("[email login] verify fail count over; email : %v ", req.Email)
			http.Err(c, consts.ErrHydraLcpUserLoginFreezed, "verify fail count over "+req.Email)
			return
		}

		// verify image captcha
		if req.ImageCaptchaID != "" {
			if err := h.phoneSrv.VerifyImageCaptcha(c, req.ImageCaptchaID, req.ImageCaptchaContent); err != nil {
				http.ErrFromGrpc(c, err)
				return
			}
		}

		userID, err = h.userSrv.VerifyPasswordByEmail(c, req.Email, req.Password)
		if err != nil {
			// 验证失败次数加1
			_, verifyFailCount := h.addVerifyFailCount(c, req.Email)
			http.Err(c, consts.ErrHydraLcpLoginFailed, fmt.Sprintf("%d", h.VerifyFailMax)+"-"+fmt.Sprintf("%d", verifyFailCount))
			return
		}
		// 清理验证失败次数
		h.deleteVerifyFailKey(c, req.Email)
	case common.EmailSMSCode:
		// 判断是否超过错误次数
		isOverMax := h.checkVerifyFailCountOverMax(c, req.Email)
		if isOverMax {
			logger.Infof("[sms code login] verify fail count over; email : %v ", req.Email)
			http.Err(c, consts.ErrHydraLcpUserLoginFreezed, "verify fail count over "+req.Email)
			return
		}

		// 验证图片验证码
		if req.ImageCaptchaID != "" {
			if err := h.phoneSrv.VerifyImageCaptcha(c, req.ImageCaptchaID, req.ImageCaptchaContent); err != nil {
				http.ErrFromGrpc(c, err)
				return
			}
		}

		// 验证手机验证码
		err := h.emailSrv.VerifyCode(c, req.Email, req.EmailCode)
		if err != nil {
			// 验证失败次数加1
			_, verifyFailCount := h.addVerifyFailCount(c, req.Email)
			http.Err(c, consts.ErrHydraLcpLoginFailed, fmt.Sprintf("%d", h.VerifyFailMax)+"-"+fmt.Sprintf("%d", verifyFailCount))
			return
		}

		// 获取用户ID
		userID, err = h.userSrv.GetID(c, req.Email, "", "")
		if err != nil {
			http.ErrFromGrpc(c, err)
			return
		}

		// 清理验证失败次数
		h.deleteVerifyFailKey(c, req.Email)
	case common.PhonePassword:
		// 判断是否超过错误次数
		isOverMax := h.checkVerifyFailCountOverMax(c, req.Phone)
		if isOverMax {
			logger.Infof("[phone password login] verify fail count over; phone : %v ", req.Phone)
			http.Err(c, consts.ErrHydraLcpUserLoginFreezed, "verify fail count over "+req.Email)
			return
		}

		// verify image captcha
		if req.ImageCaptchaID != "" {
			if err := h.phoneSrv.VerifyImageCaptcha(c, req.ImageCaptchaID, req.ImageCaptchaContent); err != nil {
				http.ErrFromGrpc(c, err)
				return
			}
		}
		// verify password
		userID, err = h.userSrv.VerifyPasswordByPhone(c, req.Phone, req.Password)
		if err != nil {
			// 手机号未找到用户，尝试使用用户名登录
			if status.Code(err) == consts.ErrHydraLcpDBUserNotExist {
				userID, err = h.userSrv.VerifyPasswordByName(c, req.Phone, req.Password)

			}
		}
		if err != nil {
			// 验证失败次数加1
			_, verifyFailCount := h.addVerifyFailCount(c, req.Phone)
			http.Err(c, consts.ErrHydraLcpLoginFailed, fmt.Sprintf("%d", h.VerifyFailMax)+"-"+fmt.Sprintf("%d", verifyFailCount))
			return
		}
		// 清理验证失败次数
		h.deleteVerifyFailKey(c, req.Phone)
	case common.PhoneSMSCode:
		// 判断是否超过错误次数
		isOverMax := h.checkVerifyFailCountOverMax(c, req.Phone)
		if isOverMax {
			logger.Infof("[sms code login] verify fail count over; phone : %v ", req.Phone)
			http.Err(c, consts.ErrHydraLcpUserLoginFreezed, "verify fail count over "+req.Phone)
			return
		}

		// 验证图片验证码
		if req.ImageCaptchaID != "" {
			if err := h.phoneSrv.VerifyImageCaptcha(c, req.ImageCaptchaID, req.ImageCaptchaContent); err != nil {
				http.ErrFromGrpc(c, err)
				return
			}
		}

		// 验证手机验证码
		err := h.phoneSrv.VerifyCode(c, req.Phone, req.PhoneCode)
		if err != nil {
			// 验证失败次数加1
			_, verifyFailCount := h.addVerifyFailCount(c, req.Phone)
			http.Err(c, consts.ErrHydraLcpLoginFailed, fmt.Sprintf("%d", h.VerifyFailMax)+"-"+fmt.Sprintf("%d", verifyFailCount))
			return
		}

		// 获取用户ID
		userID, err = h.userSrv.GetID(c, "", req.Phone, "")
		if err != nil {
			http.ErrFromGrpc(c, err)
			return
		}

		// 清理验证失败次数
		h.deleteVerifyFailKey(c, req.Phone)
	case common.Ldap:
		vv, startup := os.LookupEnv("LDAP_STARTUP")
		if startup && vv == "yes" {
			// 判断是否超过错误次数
			isOverMax := h.checkVerifyFailCountOverMax(c, req.LdapCN)
			if isOverMax {
				logger.Infof("[ldap login] verify fail count over; ldap user : %v ", req.LdapCN)
				http.Err(c, consts.ErrHydraLcpUserLoginFreezed, "verify fail count over "+req.LdapCN)
				return
			}

			ldapUserID, err := h.ldapSrv.VerifyPassword(c, req.LdapCN, req.Password)
			if err != nil {
				if status.Code(err) == consts.ErrHydraLcpLoginFailed {
					// 验证失败次数加1
					_, verifyFailCount := h.addVerifyFailCount(c, req.LdapCN)
					http.Err(c, consts.ErrHydraLcpLoginFailed, fmt.Sprintf("%d", h.VerifyFailMax)+"-"+fmt.Sprintf("%d", verifyFailCount))
					return
				}
				http.ErrFromGrpc(c, err)
				return
			}

			externalUserModel := &localModels.SsoExternalUser{UserName: req.LdapCN}
			err = service.ExternalUserService.Get(c, externalUserModel)
			if err != nil {
				http.ErrFromGrpc(c, err)
				return
			}
			// not exist, add user
			if externalUserModel.Ysid == 0 {
				id, err := h.Idgen.GenerateID(c, &idgen.GenRequest{})
				if err != nil {
					http.ErrFromGrpc(c, err)
				}
				externalUserModel.Ysid = id.Id
				externalUserModel.UserId = ldapUserID
				if err = service.ExternalUserService.Add(c, externalUserModel); err != nil {
					http.ErrFromGrpc(c, err)
					return
				}
			}
			userID = externalUserModel.Ysid

			// 清理验证失败次数
			h.deleteVerifyFailKey(c, req.LdapCN)

		} else {
			http.ErrFromGrpc(c, status.Error(consts.ErrHydraLcpLDAPNotSupported, "ldap not supported"))
			return
		}
	}

	var sid string
	sid = snowflake.ID(userID).String()
	if req.WithHydra && req.LoginType != common.Ldap {
		// sync login state to hydra here for outer applications
		logger.Infof("[api login] start api login for user %v", userID)
		res, err := h.HydraClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
			WithLoginChallenge(req.Challenge).
			WithBody(&models.HandledLoginRequest{
				Subject:     &sid,
				Remember:    true,
				RememberFor: h.TokenExpireTime,
			}))

		if err != nil {
			logger.Warnf("[api login exception] unable to send request to hydra for accepting login, err: %v", err)
			http.Errf(c, consts.ErrHydraLcpFailedToReqHydra, "unable to send request to hydra for accepting login, err: %v", err)
			return
		}

		http.Ok(c, LoginResp{RedirectURL: res.Payload.RedirectTo})
		logger.Infof("[api login] api login successful for user %v", userID)
	} else if req.WithHydra && req.LoginType == common.Ldap {
		// ldap is only used in portal applications, so sync login state to hydra-portal here
		logger.Infof("[api-portal login] start api-portal login for user %v", req.LdapCN)

		hydraCtx := make(map[string]interface{})
		hydraCtx["user_name"] = req.LdapCN
		res, err := h.HydraClientPortal.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
			WithLoginChallenge(req.Challenge).
			WithBody(&models.HandledLoginRequest{
				Subject:     &sid,
				Remember:    true,
				RememberFor: h.TokenExpireTimePortal,
				Context:     hydraCtx,
			}))

		if err != nil {
			logger.Warnf("[api-portal login exception] unable to send request to hydra-portal for accepting login, err: %v", err)
			http.Errf(c, consts.ErrHydraLcpFailedToReqHydra, "unable to send request to hydra-portal for accepting login, err: %v", err)
			return
		}

		http.Ok(c, LoginResp{RedirectURL: res.Payload.RedirectTo})
		logger.Infof("[api-portal login] api-portal login successful for user %v", userID)
	} else {
		// encode ysid with base58
		id58 := snowflake.ID(userID).String()
		http.Ok(c, id58)
	}
}

func (h *Handler) AutomaticLogin(c context.Context, req SignupAndLoginReq) (redirectURL, userId string, code int, err error) {

	var userID int64
	logger := logging.GetLogger(c)
	// keyid 存在，解密
	if req.KeyID != "" && req.Password != "" {
		cache := boot.MW.DefaultCache()
		var loginKey string
		_, ok := cache.Get(common.HydraLcpLoginKey, req.KeyID, &loginKey)
		if !ok {
			return "", "", consts.ErrHydraLcpLoginKeyExpire, errors.New("productId is a required field")
		}
		cipherPwd, err := base64.StdEncoding.DecodeString(req.Password)
		if err != nil {
			return "", "", consts.ErrHydraLcpLoginKeyExpire, errors.New("invalid base64 encoding")
		}
		pwd, err := common.AESDecrypt(cipherPwd, []byte(loginKey))

		if err != nil {
			return "", "", consts.ErrHydraLcpLoginKeyExpire, errors.New("failed to decrypt")
		}
		req.Password = string(pwd)
	}

	mtx.Lock()
	defer mtx.Unlock()

	// check username and password
	switch req.LoginType {
	case common.EmailPassword:

		userID, err = h.userSrv.VerifyPasswordByEmail(c, req.Email, req.Password)
		if err != nil {
			return "", "", consts.ErrHydraLcpLoginFailed, errors.New(fmt.Sprintf("注册失败，请重新注册！"))
		}
	case common.PhonePassword:
		// verify password
		userID, err = h.userSrv.VerifyPasswordByPhone(c, req.Phone, req.Password)
		if err != nil {
			// 手机号未找到用户，尝试使用用户名登录
			if status.Code(err) == consts.ErrHydraLcpDBUserNotExist {
				userID, err = h.userSrv.VerifyPasswordByName(c, req.Phone, req.Password)

			}
		}
		if err != nil {
			return "", "", consts.ErrHydraLcpLoginFailed, errors.New(fmt.Sprintf("注册失败，请重新注册！"))
		}
	case common.PhoneSMSCode:
		// 验证手机验证码
		err := h.phoneSrv.VerifyCode(c, req.Phone, req.Code)
		if err != nil {
			return "", "", consts.ErrHydraLcpLoginFailed, errors.New(fmt.Sprintf("注册失败，请重新注册！"))
		}

		// 获取用户ID
		userID, err = h.userSrv.GetID(c, "", req.Phone, "")
		if err != nil {
			return "", "", 0, err
		}

	case common.Ldap:
		vv, startup := os.LookupEnv("LDAP_STARTUP")
		if startup && vv == "yes" {
			// 判断是否超过错误次数
			isOverMax := h.checkVerifyFailCountOverMax(c, req.LdapCN)
			if isOverMax {
				logger.Infof("[ldap login] verify fail count over; ldap user : %v ", req.LdapCN)
				return "", "", consts.ErrHydraLcpUserLoginFreezed, errors.New("verify fail count over " + req.LdapCN)
			}

			ldapUserID, err := h.ldapSrv.VerifyPassword(c, req.LdapCN, req.Password)
			if err != nil {
				if status.Code(err) == consts.ErrHydraLcpLoginFailed {
					// 验证失败次数加1
					_, verifyFailCount := h.addVerifyFailCount(c, req.LdapCN)
					return "", "", consts.ErrHydraLcpLoginFailed, errors.New(fmt.Sprintf("%d", h.VerifyFailMax) + "-" + fmt.Sprintf("%d", verifyFailCount))
				}
				return "", "", 0, err
			}

			externalUserModel := &localModels.SsoExternalUser{UserName: req.LdapCN}
			err = service.ExternalUserService.Get(c, externalUserModel)
			if err != nil {
				return "", "", 0, err
			}
			// not exist, add user
			if externalUserModel.Ysid == 0 {
				id, err := h.Idgen.GenerateID(c, &idgen.GenRequest{})
				if err != nil {
					return "", "", 0, err
				}
				externalUserModel.Ysid = id.Id
				externalUserModel.UserId = ldapUserID
				if err = service.ExternalUserService.Add(c, externalUserModel); err != nil {
					return "", "", 0, err
				}
			}
			userID = externalUserModel.Ysid

			// 清理验证失败次数
			h.deleteVerifyFailKey(c, req.LdapCN)

		} else {
			return "", "", consts.ErrHydraLcpLDAPNotSupported, errors.New("ldap not supported")
		}
	}

	var sid string
	sid = snowflake.ID(userID).String()
	if req.WithHydra && req.LoginType != common.Ldap {
		// sync login state to hydra here for outer applications
		logger.Infof("[api login] start api login for user %v", userID)
		res, err := h.HydraClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
			WithLoginChallenge(req.Challenge).
			WithBody(&models.HandledLoginRequest{
				Subject:     &sid,
				Remember:    true,
				RememberFor: h.TokenExpireTime,
			}))

		if err != nil {
			logger.Warnf("[api login exception] unable to send request to hydra for accepting login, err: %v", err)
			return "", "", consts.ErrHydraLcpFailedToReqHydra, errors.New(fmt.Sprintf("unable to send request to hydra for accepting login, err: %v", err))
		}
		logger.Infof("[api login] api login successful for user %v", userID)
		return res.Payload.RedirectTo, "", 0, nil
	} else if req.WithHydra && req.LoginType == common.Ldap {
		// ldap is only used in portal applications, so sync login state to hydra-portal here
		logger.Infof("[api-portal login] start api-portal login for user %v", req.LdapCN)

		hydraCtx := make(map[string]interface{})
		hydraCtx["user_name"] = req.LdapCN
		res, err := h.HydraClientPortal.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
			WithLoginChallenge(req.Challenge).
			WithBody(&models.HandledLoginRequest{
				Subject:     &sid,
				Remember:    true,
				RememberFor: h.TokenExpireTimePortal,
				Context:     hydraCtx,
			}))

		if err != nil {
			logger.Warnf("[api-portal login exception] unable to send request to hydra-portal for accepting login, err: %v", err)
			return "", "", consts.ErrHydraLcpFailedToReqHydra, errors.New(fmt.Sprintf("unable to send request to hydra-portal for accepting login, err: %v", err))
		}
		logger.Infof("[api-portal login] api-portal login successful for user %v", userID)
		return res.Payload.RedirectTo, "", 0, nil
	} else {
		// encode ysid with base58
		id58 := snowflake.ID(userID).String()
		return "", id58, 0, nil
	}
}

// 判断是否超过最大验证失败次数
func (h *Handler) checkVerifyFailCountOverMax(ctx context.Context, key string) bool {
	// 判断是否超过错误次数
	cache := boot.MW.DefaultCache()
	var verifyFailCount int64
	_, ok := cache.Get(common.RedisPrefixVerifyFailPerUser, key, &verifyFailCount)
	if ok {
		if verifyFailCount >= h.VerifyFailMax {
			return true
		}
	}
	return false
}

// 增加验证失败次数
func (h *Handler) addVerifyFailCount(ctx context.Context, key string) (bool, int64) {
	cache := boot.MW.DefaultCache()
	var verifyFailCount int64
	cache.Get(common.RedisPrefixVerifyFailPerUser, key, &verifyFailCount)
	verifyFailCount++
	cache.PutWithExpire(common.RedisPrefixVerifyFailPerUser, key, verifyFailCount, time.Second*time.Duration(h.VerifyFailOverMaxFreezeLoginTime))
	return false, verifyFailCount
}

// 删除验证失败key
func (h *Handler) deleteVerifyFailKey(ctx context.Context, key string) {
	cache := boot.MW.DefaultCache()
	cache.Delete(common.RedisPrefixVerifyFailPerUser, key)
}

// LoginReqWrapper LoginReqWrapper
// swagger:parameters LoginReq
type LoginReqWrapper struct {
	// in: body
	Req LoginReq
}

// LoginReq LoginReq
type LoginReq struct {
	Challenge           string           `json:"challenge"`
	Email               string           `json:"email"`
	EmailCode           string           `json:"email_code"`
	Phone               string           `json:"phone"`
	PhoneCode           string           `json:"phone_code"`
	Password            string           `json:"password"`
	WechatUnionID       string           `json:"wechat_union_id"`
	LdapCN              string           `json:"ldap_cn"`
	LoginType           common.LoginType `json:"login_type"`
	Captcha             string           `json:"captcha"`
	CaptchaID           string           `json:"captcha_id"`
	WithHydra           bool             `json:"with_hydra"`
	ImageCaptchaID      string           `json:"image_captcha_id"`
	ImageCaptchaContent string           `json:"image_captcha_content"`
	KeyID               string           `json:"key_id"`
}

// LoginResp LoginResp
// swagger:response LoginResp
type LoginResp struct {
	RedirectURL string `json:"redirect_url"`
}

// ForwardToDefaultResp ForwardToDefaultResp
// swagger:response ForwardToDefaultResp
type ForwardToDefaultResp struct {
	DefaultURL string `json:"default_url"`
}

// WechatLoginEndpoint WechatLoginEndpoint
// swagger:route GET /api/wechat fe wechat
//
// get wechat login endpoint
//
//	 Responses:
//			200: WechatLoginEndpointResp
//
// @GET /api/wechat
func (h *Handler) WechatLoginEndpoint(c *gin.Context) {
	// "production" for production env, "dev" for dev env, "test" for test env
	// if not set, use "production"
	wechatEnv := os.Getenv("WECHAT_LOGIN_ENV")
	if wechatEnv == "" {
		wechatEnv = "production"
	}
	http.Ok(c, fmt.Sprintf(common.WechatCodeGrantURL, h.conf.Wechat.AppID, wechatEnv, url.QueryEscape(common.WechatCallbackURL)))
}

// WechatLoginEndpointRespWrapper WechatLoginEndpointRespWrapper
// swagger:response WechatLoginEndpointResp
type WechatLoginEndpointRespWrapper struct {
	// in: body
	Resp struct {
		http.SwaggerRespMeta
		// wechat login endpoint url
		Data string `json:"data"`
	}
}

// WechatLoginCallBack WechatLoginCallBack
// swagger:route GET /api/wechat/callback wechat WechatLoginCallBackRespWrapper
//
//		oauth callback api
//
//	  Responses:
//			200:
//			90019: ErrHydraLcpWechatGetOpenID
//
// @GET /api/wechat/callback
func (h *Handler) WechatLoginCallBack(c *gin.Context) {
	logger := logging.GetLogger(c)
	logger.Infof("[wechat login] user request url is %v", c.Request.RemoteAddr)

	code := c.Query("code")
	challenge := c.Query("challenge")

	url := fmt.Sprintf(common.WechatGetOpenIDURL, h.conf.Wechat.AppID, h.conf.Wechat.AppSecret, code)
	logger.Infof("[wechat login] request to wechat is %v", url)

	resp, err := http2.Get(url)
	if err != nil {
		logger.Warnf("[wechat login exception] fail to get openid: %v", err)
		http.Errf(c, consts.ErrHydraLcpWechatGetOpenID, "fail to get openid: %v", err)
		return
	}
	logger.Info("[wechat login] get openid success")

	b, err := ioutil.ReadAll(resp.Body)
	err = resp.Body.Close()
	if err != nil {
		logger.Warnf("[wechat login exception] fail to close response body: %v", err)
		http.Errf(c, consts.ErrHydraLcpWechatGetOpenID, "fail to close response body: %v", err)
		return
	}

	var r WechatTokenResp

	err = json.Unmarshal(b, &r)
	if err != nil {
		logger.Warnf("[wechat login exception] fail to unmarshal response: %v", err)
		http.Errf(c, consts.ErrHydraLcpWechatGetOpenID, "fail to unmarshal response: %v", err)
		return
	}
	logger.Infof("[wechat login] user union id is %v", r.UnionID)

	// fetch ysid by wechat_union_id
	ysid, err := h.userSrv.GetID(c, "", "", r.UnionID)
	if err != nil {
		logger.Warnf("[wechat login exception] fail to get ysid by wechat_union_id: %v", err)
		http.Errf(c, consts.ErrHydraLcpGetYSIDByUnionID, "fail to get ysid by wechat_union_id: %v", err)
		return
	}
	// get base58 encode user id
	sid := snowflake.ID(ysid).String()

	// login at hydra
	logger.Infof("[wechat login] start wechat login for user %v", ysid)
	res, err := h.HydraClient.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
		WithLoginChallenge(challenge).
		WithBody(&models.HandledLoginRequest{
			Subject:     &sid,
			Remember:    true,
			RememberFor: h.TokenExpireTime,
		}))

	if err != nil {
		logger.Warnf("[wechat login exception] unable to send request to hydra for accepting login, err: %v", err)
		http.Errf(c, consts.ErrHydraLcpFailedToReqHydra, "unable to send request to hydra for accepting login, err: %v", err)
		return
	}

	// login success and redirect
	logger.Infof("[wechat login] redirect url is %v", res.Payload.RedirectTo)
	c.Redirect(http2.StatusFound, res.Payload.RedirectTo)

	logger.Infof("[wechat login] api login successful for user %v", ysid)
}

// WechatTokenResp WechatTokenResp
type WechatTokenResp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
}

// WechatLoginCallBackResp WechatLoginCallBackResp
type WechatLoginCallBackResp struct {
	// hydra login challenge
	Challenge string `json:"challenge"`
	// union id
	// example: ohiN8wTfa1VPGaVCDuhW2e955_hM
	UnionID string `json:"unionid"`
}

// WechatLoginCallBackRespWrapper WechatLoginCallBackRespWrapper
type WechatLoginCallBackRespWrapper struct {
	// in: body
	Resp struct {
		http.SwaggerRespMeta
		Data WechatLoginCallBackResp `json:"data"`
	}
}
