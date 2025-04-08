package common

import (
	"errors"

	api "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
)

// 公共错误
var (
	// ErrInternalServer 服务器内部错误
	ErrInternalServer = errors.New(api.InternalServerErrorCode)
	// ErrInvalidArgument 无效的参数
	ErrInvalidArgument = errors.New(api.InvalidArgumentErrorCode)
	// ErrMissingArgument 缺少参数
	ErrMissingArgument = errors.New(api.MissingArgumentErrorCode)
	// ErrAccessDenied 访问拒绝
	ErrAccessDenied = errors.New(api.AccessDeniedErrorCode)
	// ErrUserNotExists 用户不存在
	ErrUserNotExists = errors.New(api.UserNotExistsErrorCode)
	// ErrAppIDNotFound appID找不到
	ErrAppIDNotFound = errors.New(api.AppIDNotFoundErrorCode)
	// ErrInvalidPath 路径不合法
	ErrInvalidPath = errors.New(api.InvalidPath)
	// ErrPathNotFound 路径文件不存在
	ErrPathNotFound = errors.New(api.PathNotFound)
	// InvalidUserID 用户ID不合法
	ErrInvalidUserID = errors.New(api.InvalidUserID)
)

// 作业相关错误
var (
	// ErrQuotaExhaustedResource 作业资源配额不足
	ErrQuotaExhaustedResource = errors.New(api.QuotaExhaustedResource)
	// ErrInvalidArgumentName 作业名称不合法
	ErrInvalidArgumentName = errors.New(api.InvalidArgumentName)
	// ErrInvalidArgumentComment 作业备注不合法
	ErrInvalidArgumentComment = errors.New(api.InvalidArgumentComment)
	// ErrInvalidArgumentResource 作业资源不合法
	ErrInvalidArgumentResource = errors.New(api.InvalidArgumentResource)
	// ErrInvalidArgumentZone 作业分区不合法
	ErrInvalidArgumentZone = errors.New(api.InvalidArgumentZone)
	// ErrInvalidArgumentInput 作业输入参数不合法
	ErrInvalidArgumentInput = errors.New(api.InvalidArgumentInput)
	// ErrInvalidArgumentOutput 作业输出参数不合法
	ErrInvalidArgumentOutput = errors.New(api.InvalidArgumentOutput)
	// ErrInvalidArgumentCommand 作业命令不合法
	ErrInvalidArgumentCommand = errors.New(api.InvalidArgumentCommand)
	// ErrInvalidArgumentEnv 作业环境变量不合法
	ErrInvalidArgumentEnv = errors.New(api.InvalidArgumentEnv)
	// ErrInvalidArgumentAllocType AllocType字段填的不对
	ErrInvalidArgumentAllocType = errors.New(api.InvalidArgumentAllocType)
	// ErrInvalidPageSize 分页大小参数不合法
	ErrInvalidPageSize = errors.New(api.InvalidPageSize)
	// ErrInvalidPageOffset 文件分片偏移量不合法
	ErrInvalidPageOffset = errors.New(api.InvalidOffset)
	// ErrJobIDNotFound 作业不存在
	ErrJobIDNotFound = errors.New(api.JobIDNotFound)
	// ErrInvalidJobID 作业ID不合法
	ErrInvalidJobID = errors.New(api.InvalidJobID)
	// ErrJobAccessDenied  作业拒绝访问
	ErrJobAccessDenied = errors.New(api.JobAccessDenied)
	// ErrJobStateNotAllowDelete 作业状态不允许删除
	ErrJobStateNotAllowDelete = errors.New(api.JobStateNotAllowDelete)
	// ErrJobStateNotAllowTerminate 作业状态不允许终止
	ErrJobStateNotAllowTerminate = errors.New(api.JobStateNotAllowTerminate)
	// ErrJobStateNotAllowResume 作业状态不允许恢复
	ErrJobStateNotAllowResume = errors.New(api.JobStateNotAllowResume)
	// ErrJobStateNotAllowQuery 作业状态不允许查询
	ErrJobStateNotAllowQuery = errors.New(api.JobStateNotAllowQuery)
	// ErrJobStateNotAllowTransmitSuspend 作业状态不允许传输暂停
	ErrJobStateNotAllowTransmitSuspend = errors.New(api.JobStateNotAllowTransmitSuspend)
	// ErrJobStateNotAllowTransmitResume 作业状态不允许传输恢复
	ErrJobStateNotAllowTransmitResume = errors.New(api.JobStateNotAllowTransmitResume)
	// ErrJobPathUnauthorized 目标路径未授权
	ErrJobPathUnauthorized = errors.New(api.JobPathUnauthorized)
	// ErrJobFileSyncStateUpdateFailed 更新作业文件传输状态异常
	ErrJobFileSyncStateUpdateFailed = errors.New(api.JobFileSyncStateUpdateFailed)
	// ErrJobNotAllowedRetransmit 作业文件传输不允许重新开始
	ErrJobNotAllowedRetransmit = errors.New(api.JobNotAllowedRetransmit)
	// ErrInvalidArgumentDownloadFinishedTime 下载完成时间格式异常
	ErrInvalidArgumentDownloadFinishedTime = errors.New(api.InvalidArgumentDownloadFinishedTime)
	// ErrInvalidArgumentDownloadFileSizeCurrent 下载文件大小错误
	ErrInvalidArgumentDownloadFileSizeCurrent = errors.New(api.InvalidArgumentDownloadFileSizeTotal)
	// ErrInvalidArgumentDownloadFileSizeTotal 下载文件总大小错误
	ErrInvalidArgumentDownloadFileSizeTotal = errors.New(api.InvalidArgumentDownloadFileSizeTotal)
	// ErrInvalidArgumentCustomStateRuleKeyStatement 作业自定义状态规则key语句不合法
	ErrInvalidArgumentCustomStateRuleKeyStatement = errors.New(api.InvalidArgumentCustomStateRuleKeyStatement)
	// ErrInvalidArgumentCustomStateRuleResultState 作业自定义状态规则resultState不合法
	ErrInvalidArgumentCustomStateRuleResultState = errors.New(api.InvalidArgumentCustomStateRuleResultState)
	// ErrInvalidArgumentJobIDs 作业IDs不合法
	ErrInvalidArgumentJobIDs = errors.New(api.InvalidArgumentJobIDs)
	// ErrInvalidArgumentPageSize 作业分页大小不合法
	ErrInvalidArgumentPageSize = errors.New(api.InvalidPageSize)
	// ErrInvalidArgumentPageOffset 作业分页偏移量不合法
	ErrInvalidArgumentPageOffset = errors.New(api.InvalidPageOffset)
	// ErrInvalidArgumentJobState 作业状态不合法
	ErrInvalidArgumentJobState = errors.New(api.InvalidArgumentJobState)
	// ErrInvalidPreScheduleID 预调度ID不合法
	ErrInvalidPreScheduleID = errors.New(api.InvalidPreScheduleID)
	// ErrPreScheduleNotFound 预调度不存在
	ErrPreScheduleNotFound = errors.New(api.PreScheduleNotFound)
	// ErrPreScheduleUsed 预调度已使用
	ErrPreScheduleUsed = errors.New(api.PreScheduleUsed)
	// ErrWrongCPUUsage CPU使用率获取错误
	ErrWrongCPUUsage = errors.New(api.WrongCPUUsage)

	// for application
	ErrInvalidAppID   = errors.New(api.InvalidAppID)
	ErrDuplicateEntry = errors.New("app version,type already exists")
	ErrNoEffect       = errors.New("no effect")

	// ErrAppQuotaNotFound 应用配额不存在
	ErrAppQuotaNotFound = errors.New(api.AppQuotaNotFound)
	// ErrAppQuotaAlreadyExist 应用配额已存在
	ErrAppQuotaAlreadyExist = errors.New(api.AppQuotaAlreadyExist)
	// ErrUserNoAppQuota 用户没有应用配额
	ErrUserNoAppQuota = errors.New(api.UserNoAppQuota)
	// ErrAppNotPublished 应用未发布
	ErrAppNotPublished         = errors.New(api.AppNotPublished)
	ErrJobResidualNotFound     = errors.New(api.JobResidualNotFound)
	ErrHpcResidual             = errors.New(api.JobGetHpcResidualFailed)
	ErrJobSnapshotNotFound     = errors.New(api.JobSnapshotNotFound)
	ErrJobMonitorChartNotFound = errors.New(api.JobMonitorChartNotFound)

	// ErrAppAllowNotFound 应用白名单不存在
	ErrAppAllowNotFound = errors.New(api.AppAllowNotFound)
	// AppAllowAlreadyExist 应用白名单已存在
	ErrAppAllowAlreadyExist = errors.New(api.AppAllowAlreadyExist)
)

var (
	ErrInvalidChargeParams                  = errors.New(api.InvalidArgumentChargeParams)
	ErrInvalidAccountId                     = errors.New(api.InvalidArgumentAccountId)
	ErrInvalidAccountStatusNotEnoughBalance = errors.New(api.InvalidAccountStatusNotEnoughBalance)
	ErrInvalidAccountStatusFrozen           = errors.New(api.InvalidAccountStatusFrozen)
)

var (
	ErrInvalidPayBy                  = errors.New(api.InvalidArgumentPayBy)
	ErrPayByTokenExpire              = errors.New(api.PayByTokenExpire)
	ErrInvalidArgumentPayBySignature = errors.New(api.InvalidArgumentPayBySignature)
)

var (
	ErrLicenseIDNotFound = errors.New(api.LicenseIdNotFound)
)

// 账号相关错误
var (
	// ErrAccountNotExists 账户不存在
	ErrAccountNotExists = errors.New(consts.AccountNotExists)
	// ErrAccountBillNotExists 账单不存在
	ErrAccountBillNotExists = errors.New(consts.AccountBillNotExists)
	// ErrAccountAmount  金额不合法
	ErrAccountAmount = errors.New(consts.InvalidAmount)
	// ErrInvalidComment 备注不合法
	ErrInvalidComment = errors.New(consts.InvalidComment)
	// ErrInDebt 账号已欠费
	ErrInDebt = errors.New(consts.InDebt)
	// ErrAccountExists 账号已存在
	ErrAccountExists = errors.New(consts.AccountExists)
	// ErrCreditAddTradeExists 充值操作已存在
	ErrCreditAddTradeExists = errors.New(consts.CreditAddTradeExists)
	// ErrCreditAddTradeExists 扣减操作已存在
	ErrReduceTradeExists = errors.New(consts.ReduceTradeExists)
	// ErrCreditQuotaExhausted 授信额度超出
	ErrCreditQuotaExhausted = errors.New(consts.CreditQuotaExhausted)
	// ErrInsufficientBalance 账户余额不足
	ErrInsufficientBalance = errors.New(consts.InsufficientBalance)
	// ErrFreezedAccount 账户已被冻结
	ErrFreezedAccount = errors.New(consts.FreezedAccount)
	// ErrAccountBillSignStatusInvalid 账户账单操作状态异常
	ErrAccountBillSignStatusInvalid = errors.New(consts.AccountBillSignStatusInvalid)
	// ErrAccountBillIdempotentIDRepeat AccountBillIdempotentIDRepeat
	ErrAccountBillIdempotentIDRepeat = errors.New(consts.AccountBillIdempotentIDRepeat)

	ErrAccountResponse      = errors.New("invalid account response")
	ErrRequestAccountByYSID = errors.New("request account by ysid failed")
)

// 代金券相关错误
var (
	// ErrCashVoucherNotExists
	ErrCashVoucherNotExists     = errors.New(consts.CashVoucherNotExists)
	ErrAccountVoucherIDNotFound = errors.New(consts.AccountVoucherIDNotFound)
	ErrAccountVoucherExpired    = errors.New(consts.AccountVoucherExpired)
	ErrAccountVoucherDisabled   = errors.New(consts.AccountVoucherDisabled)
	ErrAccountVoucherExceed     = errors.New(consts.AccountVoucherExceed)
)

// sync-runner 相关错误
var (
	// ErrSyncTaskNotFound ...
	ErrSyncTaskNotFound = errors.New(api.SyncTaskNotFound)
)

var (
	ErrPageOffset         = errors.New("the PageOffset parameter should >= 0")
	ErrPageSize           = errors.New("the PageOffset parameter should 1 ~ 1000")
	ErrTooManyJobNum      = errors.New("user submit job nums is too many")
	ErrNoPermissionForJob = errors.New("no permission for job")
	ErrJobIsFinalState    = errors.New("job is final state")
	ErrJobIsNotFinalState = errors.New("job is not final state")
	ErrUserQuotaNotFound  = errors.New("user quota not found")
	ErrInstanceExecScript = errors.New("instance executing script failed")
)

// cad converter 相关错误
var (
	ErrNotEnoughTokensAvailable = errors.New("not enough tokens available")
	ErrAcquireTokenTimeout      = errors.New("acquire token timeout")
)
