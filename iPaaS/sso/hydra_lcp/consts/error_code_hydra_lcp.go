package consts

// error code for hydra integrated login consent provider

// from 90001 to 100000
const (
	ErrHydraLcpFailedToReqHydra = 90001
	ErrHydraLcpBadRequest       = 90002

	ErrHydraLcpLackYsID            = 90003
	ErrHydraLcpPwdNotMatch         = 90004
	ErrHydraLcpPwdFailToGenerate   = 90005
	ErrHydraLcpUtf8ToGbk           = 90006
	ErrHydraLcpJSONMarshal         = 90007
	ErrHydraLcpJSONUnmarshal       = 90008
	ErrHydraLcpSendHTTPPostReq     = 90009
	ErrHydraLcpFailToSendSMS       = 90010
	ErrHydraLcpCaptchaFailed       = 90011
	ErrHydraLcpRedisFailed         = 90012
	ErrHydraLcpWait                = 90013
	ErrHydraLcpFailedToParseTime   = 90014
	ErrHydraLcpSMTPError           = 90015
	ErrHydraLcpJWTInvalid          = 90016
	ErrHydraLcpJWTGenerate         = 90017
	ErrHydraLcpUserNotMatch        = 90018
	ErrHydraLcpWechatGetOpenID     = 90019
	ErrHydraLcpLdapError           = 90020
	ErrHydraLcpUserExist           = 90021
	ErrHydraLcpDecryptYsid         = 90022
	ErrHydraLcpNameEmpty           = 90023
	ErrHydraLcpPhoneEmpty          = 90024
	ErrHydraLcpPasswordEmpty       = 90025
	ErrHydraLcpYsidEmpty           = 90026
	ErrHydraLcpPhoneExist          = 90027
	ErrHydraLcpEmailExist          = 90028
	ErrHydraLcpCaptchaVerifyFailed = 90029
	ErrHydraLcpSignupFollowUp      = 90030
	ErrHydraLcpLackWechatInfo      = 90031
	ErrHydraLcpWechatInfoExist     = 90032
	ErrHydraLcpGetYSIDByUnionID    = 90033
	ErrHydraLcpLDAPNotSupported    = 90034
	ErrHydraLcpPhoneNotExist       = 90035
	ErrHydraLcpPwdInvalidate       = 90036 // 密码格式错误
	ErrHydraLcpChallengeExpired    = 90037 // challenge 过期
	ErrHydraLcpUserLoginFreezed    = 90038 // 用户登录锁定
	ErrHydraLcpPhoneInvalidate     = 90039 // 手机号格式错误
	ErrHydraLcpLoginFailed         = 90040 // 登录失败（用户名或密码错误）
	ErrHydraLcpLoginKeyExpire      = 90041 // 登录key过期
	ErrHydraLcpGetLoginKeyFail     = 90042 // 获取登录key失败
	ErrHydraLcpEmailNotExist       = 90046

	ErrHydraLcpIDMUserInfoFail    = 90043 // idm 获取用户信息失败
	ErrHydraLcpIDMUserPhoneNotSet = 90044 // idm 用户手机号码为空
	ErrHydraLcpIDMMqEventError    = 90045 // 发送 idm 到 mq 失败

	// database error
	ErrHydraLcpDBOpFail          = 90100
	ErrHydraLcpDBDuplicatedEntry = 90101
	ErrHydraLcpDBUserNotExist    = 90103

	// error for external user
	ErrHydraLcpExternalUserNotFound = 90201

	// error for hybrid user
	ErrHydraLcpPeerNotFound   = 90301
	ErrHydraLcpPeerInvalid    = 90302
	ErrHydraLcpSecretIsExpire = 90303

	// 短信模板不存在
	ErrHydraLcpSendSmsTplNotFound = 90321
	// 短信参数与短信模板要求不符
	ErrHydraLcpSendSmsParamCount = 90322

	// error for jwt auth
	ErrHydraLcpJwtAuthError      = 90351
	ErrHydraLcpPhoneJwtAuthError = 90352
)
