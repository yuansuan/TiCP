// This code is generated, DO NOT EDIT.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/handler"
)

// UseRoutersGenerated UseRoutersGenerated
func UseRoutersGenerated(server *gin.Engine) {
	Handler := handler.CreateHandler()
	server.POST("/api/mkevent/login", Handler.MarketEventLogin)
	server.GET("/api/mkevent/getuserinfo", Handler.GetUserInfo)
	server.GET("/api/hydra/login", Handler.HydraLogin)
	server.GET("/api/hydra_portal/login", Handler.HydraPortalLogin)
	server.GET("/api/send_code", Handler.SendVerificationCode)
	server.GET("/api/send_code_new", Handler.NewSendVerificationCode)
	server.GET("/api/forward_to_default", Handler.ForwardToDefault)
	server.POST("/api/login", Handler.Login)
	server.GET("/api/wechat", Handler.WechatLoginEndpoint)
	server.GET("/api/wechat/callback", Handler.WechatLoginCallBack)
	server.GET("/api/login/getkey", Handler.GetLoginKey)
	server.POST("/api/signup", Handler.Signup)
	server.POST("/api/signupandlogin", Handler.SignupAndLogin)
	server.POST("/api/signupbyname", Handler.SignupByName)
	server.POST("/api/resend", Handler.Resend)
	server.POST("/api/activate", Handler.Activate)
	server.POST("/api/captcha", Handler.CreateCaptcha)
	server.POST("/api/reset_password", Handler.ResetPassword)
	server.POST("/api/reset_email_password", Handler.ResetEmailPassword)
	server.GET("/api/hydra/challengetologinurl", Handler.HydraToLoginURL)
	server.GET("/api/hydra_portal/challengetologinurl", Handler.HydraPortalToLoginURL)
	server.GET("/api/hydra/consent", Handler.HydraConsent)
	server.GET("/api/hydra_portal/consent", Handler.HydraPortalConsent)
	server.GET("/api/hydra/logout", Handler.HydraLogout)
	server.GET("/api/hydra_portal/logout", Handler.HydraPortalLogout)
	server.GET("/api/offiaccount/callback", Handler.OffiaccountCallback)
	server.POST("/api/offiaccount/callback", Handler.OffiaccountCallback)
	server.POST("/api/offiaccount/createqrcode", Handler.CreateQRCodeWithParam)
	server.GET("/api/casi/login", Handler.CallCASI)
	server.GET("/api/casi/callback", Handler.CASICallback)
}

// This code is generated, DO NOT EDIT.
