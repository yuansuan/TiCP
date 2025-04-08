package common

// defined by hydra
const (
	// request from hydra
	HydraLoginChallenge   = "login_challenge"
	HydraLoginJwt         = "jwt"
	HydraConsentChallenge = "consent_challenge"
	HydraLogoutChallenge  = "logout_challenge"
)

// defined by oauth protocol
const (
	Wechat = "wechat"
	// fmt.Sprintf(WechatCodeGrantURL, app_id, state)
	WechatCodeGrantURL = `https://open.weixin.qq.com/connect/qrconnect?appid=%v&response_type=code&scope=snsapi_login&state=%v&redirect_uri=%v`
	// fmt.Sprintf(WechatGetOpenIDURL, app_id, app_secret, code)
	WechatGetOpenIDURL = `https://api.weixin.qq.com/sns/oauth2/access_token?appid=%v&secret=%v&code=%v&grant_type=authorization_code`
	WechatCallbackURL  = "http://yuansuan.cloud/api/wechat/callback"
)

const (
	// CachePrefix CachePrefix
	CachePrefix = "HydraLcpSecret"
	// RedisPrefixCodeIDToCaptcha RedisPrefixCodeIDToCaptcha
	RedisPrefixCodeIDToCaptcha = "HydraLcpCodeIDToCaptcha"
	// RedisPrefixPhoneToCodeID RedisPrefixPhoneToCodeID
	RedisPrefixPhoneToCodeID = "HydraLcpPhoneToCodeID"
	// RedisPrefixEmailToCodeID RedisPrefixEmailToCodeID
	RedisPrefixEmailToCodeID = "HydraLcpEmailToCodeID"
	// RedisPrefixPhoneToLastSendTime RedisPrefixPhoneToLastSendTime
	RedisPrefixPhoneToLastSendTime = "HydraLcpPhoneToLastSendTime"
	//RedisPrefixPhoneToLastSendTime RedisPrefixPhoneToLastSendTime
	RedisPrefixEmailToLastSendTime = "HydraLcpEmailToLastSendTime"
	// RedisPrefixImageCaptchaIDToIndex RedisPrefixImageCaptchaIDToIndex
	RedisPrefixImageCaptchaIDToIndex = "HydraLcpImageCaptchaIDToIndex"
	// RedisPrefixSendCountPerHourAndIP RedisPrefixSendCountPerHourAndIP
	RedisPrefixSendCountPerHourAndIP = "HydraLcpSendCountPerHourAndIP"
	// RedisPrefixVerifyFailPerUser RedisPrefixVerifyFailPerUser
	RedisPrefixVerifyFailPerUser = "HydraLcpVerifyFailPerUser"
	// HydraLcpLoginKey HydraLcpLoginKey
	HydraLcpLoginKey = "HydraLcpLoginKey"
)

// sms template
const (
	// the expired time is set in redis
	SMStemplateGetCode       = `您的验证码是%v，在5分钟内有效。如非本人操作请忽略本短信。`
	SMSTemplatePhoneCode     = "PhoneCode"
	SMSTemplateSignupSuccess = "SignupSuccess"
	SMSTemplateEmailCode     = "EmailCode"
	//
	// nation code
	PhoneCodeChina = "+86"
)

// LdapBaseDN LdapBaseDN
var LdapBaseDN = "ou=People,dc=yuansuan,dc=cn"

// LdapID LdapID
var LdapID = "uidNumber"

// LoginType LoginType
type LoginType int

const (
	// EmailPassword EmailPassword
	EmailPassword LoginType = iota + 1
	// PhonePassword PhonePassword
	PhonePassword
	// PhoneSMSCode PhoneSMSCode
	PhoneSMSCode
	// WeChat WeChat
	WeChat
	// Ldap Ldap
	Ldap
	// EmailSMSCode EmailSMSCode
	EmailSMSCode
)

// String String
func (m LoginType) String() string {
	LoginTypes := []string{"Email with password", "Phone with password", "Phone with SMS code", "WeChat API"}
	if m > WeChat || m < EmailPassword {
		return "Unknown"
	}
	return "Login in " + LoginTypes[m]
}

// JWTSecret JWTSecret
var JWTSecret = []byte("secret")

// JWTSignup ...
var JWTSignup = "signup"

// monitor
const (
	HydraLcpSmsSend = "hydra_lcp_sms_send"
	SendResult      = "send_result"
	ServiceProvider = "service_provider"
)

// product check
const (
	HydraLcpProductID = "product_id"
)
const (
	TypeEmail = "email"
	TypePhone = "phone"
)
