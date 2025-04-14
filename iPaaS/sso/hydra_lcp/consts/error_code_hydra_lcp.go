package consts

// error code for hydra integrated login consent provider

// from 90001 to 100000
const (
	ErrHydraLcpBadRequest        = 90002
	ErrHydraLcpLackYsID          = 90003
	ErrHydraLcpPwdNotMatch       = 90004
	ErrHydraLcpPwdFailToGenerate = 90005
	ErrHydraLcpUserExist         = 90021
	ErrHydraLcpNameEmpty         = 90023
	ErrHydraLcpPhoneEmpty        = 90024
	ErrHydraLcpPasswordEmpty     = 90025
	ErrHydraLcpYsidEmpty         = 90026
	ErrHydraLcpLackWechatInfo    = 90031
	ErrHydraLcpWechatInfoExist   = 90032
	ErrHydraLcpPwdInvalidate     = 90036 // 密码格式错误
	ErrHydraLcpPhoneInvalidate   = 90039 // 手机号格式错误
	// database error
	ErrHydraLcpDBOpFail          = 90100
	ErrHydraLcpDBDuplicatedEntry = 90101
	ErrHydraLcpDBUserNotExist    = 90103
)
