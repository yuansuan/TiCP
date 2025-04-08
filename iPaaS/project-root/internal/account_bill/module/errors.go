package module

import "errors"

// module 错误定义
var (
	ErrServerInternal = errors.New("内部错误")

	ErrAccountNotExist   = errors.New("账户不存在")
	ErrAccountNotUpdated = errors.New("账户没有信息更新或授信额度小于已透支金额")
	ErrAccountAdd        = errors.New("新增账号异常")
)

var (
	ErrCashVoucherIDNotFound    = errors.New("代金券不存在")
	ErrAccountVoucherIDNotFound = errors.New("账户代金券不存在")
)
