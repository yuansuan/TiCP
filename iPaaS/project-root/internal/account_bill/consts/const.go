package consts

// 公共错误码
const (
	// InternalServerErrorCode 服务器内部错误
	InternalServerErrorCode = "InternalServerError"
	// InvalidArgumentErrorCode  无效的参数
	InvalidArgumentErrorCode = "InvalidArgument"
	// MissingArgumentErrorCode 缺少参数
	MissingArgumentErrorCode = "MissingArgument"
	// AccessDeniedErrorCode 访问拒绝
	AccessDeniedErrorCode = "AccessDenied"
	// UserNotExistsErrorCode 用户不存在
	UserNotExistsErrorCode = "UserNotExists"
	// AppIDNotFoundErrorCode appID找不到
	AppIDNotFoundErrorCode = "AppIdNotFound"
	// InvalidPageSizeErrorCode 分页大小参数不合法
	InvalidPageSizeErrorCode = "InvalidArgument.PageSize"
	// InvalidPageIndexErrorCode 分页偏移量不合法
	InvalidPageIndexErrorCode = "InvalidArgument.PageIndex"
	// InvalidUserId 非法的用户ID
	InvalidUserId = "Invalid.UserId"
	// HydraLcpDBUserNotExist 用户不存在
	ErrHydraLcpDBUserNotExist = 90103
)

const (
	ACCONT_ID_KEY                    = "AccountID"
	ACCOUNT_USER_ID_KEY              = "UserID"
	Account_CashVoucher_ID_KEY       = "AccountCashVoucherID"
	AccountBillUniqueIndexIdempotent = "account_bill__uindex_idempotent"
)

// 账户错误码
const (
	// AccountNotExists  账号不存在
	AccountNotExists = "AccountNotExists"
	// AccountBillNotExists 账单不存在
	AccountBillNotExists = "AccountBillNotExists"
	// InvalidAmount 金额不能小于0
	InvalidAmount = "InvalidArgument.Amount"
	// InvalidComment 不能为空
	InvalidComment = "InvalidArgument.comment"
	// InvalidEndTime 结束时间不能小于开始时间
	InvalidEndTime = "InvalidArgument.EndTime"
	// InDebt 账户已欠费
	InDebt = "InDebt"
	// AccountExists 账号已存在
	AccountExists = "AccountExists"
	// CreditAddTradeExists 充值操作已经存在
	CreditAddTradeExists = "CreditAddTradeExists"
	// ReduceTradeExists 扣减已经存在
	ReduceTradeExists = "ReduceTradeExists"
	// CreditQuotaExhausted 授信额度超出
	CreditQuotaExhausted = "CreditQuotaExhausted"
	// InsufficientBalance 余额不足
	InsufficientBalance = "InsufficientBalance"
	// FreezedAccount 账户已被冻结
	FreezedAccount = "FreezedAccount"
	// AccountBillSignStatusInvalid 账户账单操作状态异常
	AccountBillSignStatusInvalid = "AccountBillSignStatusInvalid"

	// AccountBillIdempotentIDRepeat 幂等值不唯一
	AccountBillIdempotentIDRepeat = "AccountBillIdempotentIDRepeat"
)
const (
	Success = "success"
)

const (
	AmountZero = 0
)

const (
	Status        = 1
	InvalidStatus = -1
)
const (
	Enabled = 1
)
const (
	NoExpired = 0
	Expired   = 1
)
const (
	NoDeleted = 0
	Deleted   = 1
)

const (
	AbsExpired = 1 // 1:绝对
	RelExpired = 2 // 2:相对
)

const (
	EXPIRED_COMMENT = "账户代金卷过期"
)

// 代金券
const (
	CASH_VOUCHER_ID            = "CashVoucherID"
	CASH_VOUCHER_EXCEED_AMOUNT = 10
)

const (
	VOUCHER_LOG_SIGN_CONSUME = 1
	VOUCHER_LOG_SIGN_EXPIRED = 2
)

// 单位
const (
	UNIT_PRICE_DEFAULT = "元"
)

// 代金券异常
const (
	CashVoucherNotExists     = "CashVoucherNotExists"
	AccountVoucherIDNotFound = "AccountVoucherIDNotFound"
	AccountVoucherExpired    = "AccountVoucherExpired"
	AccountVoucherDisabled   = "AccountVoucherDisabled"
	AccountVoucherExceed     = "AccountVoucherExceed"
)
