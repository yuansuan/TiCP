package common

// 公共错误码
const (
	// InternalServerErrorCode 服务器内部错误
	InternalServerErrorCode = "InternalServerError"
	// InvalidArgumentErrorCode  无效的参数
	InvalidArgumentErrorCode = "InvalidArgument"
	// MissingArgumentErrorCode 缺少参数
	MissingArgumentErrorCode = "MissingArgument"
	// InvalidParamErrorCode  无效的参数
	InvalidParamErrorCode = "InvalidParams"
	// MissingParamErrorCode 缺少参数
	MissingParamErrorCode = "MissingParams"
	// AccessDeniedErrorCode 访问拒绝
	AccessDeniedErrorCode = "AccessDenied"
	// UserNotExistsErrorCode 用户不存在
	UserNotExistsErrorCode = "UserNotExists"
	// ProjectIDNotFoundErrorCode ProjectID找不到
	ProjectIDNotFoundErrorCode = "ProjectIdNotFound"
	// ZoneNotFoundErrorCode zone找不到
	ZoneNotFoundErrorCode = "ZoneNotFound"
	// AppIDNotFoundErrorCode appID找不到
	AppIDNotFoundErrorCode = "AppIdNotFound"
	// InvalidBucketErrorCode 非法的bucket
	InvalidBucketErrorCode = "InvalidBucket"
	// InvalidPageSize 分页大小参数不合法
	InvalidPageSize = "InvalidArgument.PageSize"
	// InvalidPageOffset 分页偏移量不合法
	InvalidPageOffset = "InvalidArgument.PageOffset"
	// InvalidUserID 用户ID不合法
	InvalidUserID = "InvalidArgument.UserID"
	// InvalidBeginTime 开始时间不合法
	InvalidBeginTime = "InvalidArgument.BeginTime"
	// InvalidEndTime 结束时间不合法
	InvalidEndTime = "InvalidArgument.EndTime"

	NoEffect = "NoEffect"
)

// SystemParams 系统参数(使用|分隔)
const SystemParams = "MAIN_FILE|MAIN_FILE_CONTENT"

// 作业相关错误码
const (
	QuotaExhaustedResource                     = "QuotaExhausted.Resource"
	InvalidArgumentName                        = "InvalidArgument.Name"
	InvalidArgumentComment                     = "InvalidArgument.Comment"
	InvalidArgumentZone                        = "InvalidArgument.Zone"
	InvalidArgumentInput                       = "InvalidArgument.Input"
	InvalidArgumentOutput                      = "InvalidArgument.Output"
	InvalidArgumentCommand                     = "InvalidArgument.Command"
	InvalidArgumentEnv                         = "InvalidArgument.Env"
	InvalidArgumentJobState                    = "InvalidArgument.JobState"
	JobIDNotFound                              = "JobIdNotFound"
	JobAccessDenied                            = "JobAccessDenied"
	JobStateNotAllowDelete                     = "JobStateNotAllowDelete"
	JobStateNotAllowTerminate                  = "JobStateNotAllowTerminate"
	JobStateNotAllowResume                     = "JobStateNotAllowResume"
	JobStateNotAllowQuery                      = "JobStateNotAllowQuery"
	JobStateNotAllowTransmitSuspend            = "JobStateNotAllowTransmitSuspend"
	JobStateNotAllowTransmitResume             = "JobStateNotAllowTransmitResume"
	JobPathUnauthorized                        = "JobPathUnauthorized"
	JobResidualNotFound                        = "JobResidualNotFound"
	JobGetHpcResidualFailed                    = "JobGetHpcResidualFailed"
	JobSnapshotNotFound                        = "JobSnapshotNotFound"
	JobMonitorChartNotFound                    = "JobMonitorChartNotFound"
	JobFileSyncStateUpdateFailed               = "JobFileSyncStateUpdateFailed"
	JobNotAllowedRetransmit                    = "JobNotAllowedRetransmit"
	InvalidArgumentDownloadFinishedTime        = "InvalidArgument.DownloadFinishedTime"
	InvalidArgumentDownloadFileSizeTotal       = "InvalidArgument.DownloadFileSizeTotal"
	InvalidArgumentDownloadFileSizeCurrent     = "InvalidArgument.DownloadFileSizeCurrent"
	InvalidArgumentCustomStateRuleKeyStatement = "InvalidArgument.CustomStateRule.KeyStatement"
	InvalidArgumentCustomStateRuleResultState  = "InvalidArgument.CustomStateRule.ResultState"
	InvalidArgumentDescription                 = "InvalidArgument.Description"
	InvalidArgumentVersion                     = "InvalidArgument.Version"
	InvalidArgumentType                        = "InvalidArgument.Type"
	InvalidArgumentJobIDs                      = "InvalidArgument.JobIDs"
	InvalidArgumentResource                    = "InvalidArgument.Resource"
	InvalidArgumentAllocType                   = "InvalidArgument.AllocType"
	InvalidPreScheduleID                       = "InvalidArgument.PreScheduleID"
	PreScheduleNotFound                        = "PreScheduleNotFound"
	PreScheduleUsed                            = "PreScheduleUsed"
	WrongCPUUsage                              = "WrongCPUUsage"

	InvalidAppID                      = "InvalidArgument.AppID"
	InvalidArgumentBinPath            = "InvalidArgument.BinPath"
	InvalidArgumentExtentionParams    = "InvalidArgument.ExtentionParams"
	InvalidArgumentPublishStatus      = "InvalidArgument.PublishStatus"
	InvalidArgumentLicManagerId       = "InvalidArgument.LicManagerId"
	InvalidJobID                      = "InvalidArgument.JobID"
	InvalidArgumentResidualLogParser  = "InvalidArgument.ResidualLogParser"
	InvalidArgumentMonitorChartParser = "InvalidArgument.MonitorChartParser"

	InvalidArgumentSpecifyQueue = "InvalidArgument.SpecifyQueue"

	AppQuotaNotFound     = "AppQuotaNotFound"
	AppQuotaAlreadyExist = "AppQuotaAlreadyExist"
	UserNoAppQuota       = "UserNoAppQuota"
	AppNotPublished      = "AppNotPublished"

	AppAllowNotFound     = "AppAllowNotFound"
	AppAllowAlreadyExist = "AppAllowAlreadyExist"
)

// 存储相关错误码
const (
	// InvalidStorageVersion 存储版本不合法
	InvalidStorageVersion = "InvalidArgument.StorageVersion"
	// InvalidPath 路径参数不合法
	InvalidPath = "InvalidArgument.Path"
	// InvalidBasePath 基准路径参数不合法
	InvalidBasePath = "InvalidArgument.BasePath"
	// SizeTooLarge 文件大小超过限制
	SizeTooLarge = "InvalidArgument.SizeTooLarge"
	//LengthTooLarge 文件分片长度太大
	LengthTooLarge = "InvalidArgument.LengthTooLarge"
	// InvalidOffset 文件分片偏移量不合法
	InvalidOffset = "InvalidArgument.Offset"
	// InvalidLength 加上Length后，超过文件总长度
	InvalidLength = "InvalidArgument.Length"
	//InvalidRange Range参数不合法
	InvalidRange = "InvalidArgument.Range"
	//InvalidData 文件Data参数不合法
	InvalidData = "InvalidArgument.Data"
	//InvalidSize 文件Size参数不合法
	InvalidSize = "InvalidArgument.Size"
	//InvalidCompressor 压缩器参数不合法
	InvalidCompressor = "InvalidArgument.Compressor"
	// InvalidBlockSize 分块大小参数不合法
	InvalidBlockSize = "InvalidArgument.BlockSize"
	// InvalidBeginChunkOffset 起始分块偏移量参数不合法
	InvalidBeginChunkOffset = "InvalidArgument.BeginChunkOffset"
	// InvalidEndChunkOffset 结束分块偏移量参数不合法
	InvalidEndChunkOffset = "InvalidArgument.EndChunkOffset"
	// InvalidRollingHashType 弱校验算法类型参数不合法
	InvalidRollingHashType = "InvalidArgument.RollingHashType"
	// InvalidSrcOffset src文件分片偏移量不合法
	InvalidSrcOffset = "InvalidArgument.SrcOffset"
	// InvalidDestOffset dest文件分片偏移量不合法
	InvalidDestOffset = "InvalidArgument.DestOffset"
	//InvalidRegexp 正则表达式不合法
	InvalidRegexp = "InvalidArgument.Regexp"
	//InvalidCompressID 压缩任务ID不合法
	InvalidCompressID = "InvalidArgument.CompressID"
	//InvalidFileName 文件名不合法
	InvalidFileName = "InvalidArgument.FileName"
	//InvalidTargetPath 目标路径不合法
	InvalidTargetPath = "InvalidArgument.TargetPath"
	//InvalidStorageLimit 存储上限不合法
	InvalidStorageLimit = "InvalidArgument.StorageLimit"
	//InvalidFileType 文件类型不合法
	InvalidFileType = "InvalidArgument.FileType"
	//InvalidOperationType 操作类型不合法
	InvalidOperationType = "InvalidArgument.OperationType"
	//InvalidSyncTaskStopMode 同步任务停止模式不合法
	InvalidSyncTaskStopMode = "InvalidArgument.StopMode"
	// DestPathNotFound 目标文件/文件夹路径不存在
	DestPathNotFound = "DestPathNotFound"
	//PathNotFound 文件夹路径不存在
	PathNotFound = "PathNotFound"
	// PathExists 文件夹/文件已存在
	PathExists = "PathExists"
	// PathContainsFile 路径包含文件
	PathContainsFile = "PathContainsFile"
	// TargetPathExists 目标文件/文件夹已存在
	TargetPathExists = "TargetPathExists"
	// SrcPathNotFound 源文件/文件夹路径不存在
	SrcPathNotFound = "SrcPathNotFound"
	//DestPathExists 目标文件/文件夹路径已存在
	DestPathExists = "DestPathExists"
	//SystemQuotaExhausted 系统存储空间超过配额
	SystemQuotaExhausted = "SystemQuotaExhausted"
	//QuotaExhausted 用户存储空间超过配额
	QuotaExhausted = "QuotaExhausted"
	//StorageQuotaNotFound 存储配额没有找到
	StorageQuotaNotFound = "StorageQuotaNotFound"
	//UploadIDNotFound uploadID没有找到
	UploadIDNotFound = "UploadIDNotFound"
	//PathNotMatchUploadInit 文件路径和init时传的不一致
	PathNotMatchUploadInit = "PathNotMatchUploadInit"
	//CompressTaskExists 压缩任务已存在(一个用户同一时间只能有一个压缩任务)
	CompressTaskExists = "CompressTaskExists"
	//TooManyCompressTask 压缩任务太多
	TooManyCompressTask = "TooManyCompressTask"
	//UploadTaskExists 上传任务已存在
	UploadTaskExists = "UploadTaskExists"
	//CompressTaskNotFound 压缩任务不存在
	CompressTaskNotFound = "CompressTaskNotFound"
	//CompressTaskIsFinished 压缩任务已结束
	CompressTaskIsFinished = "CompressTaskIsFinished"
	//CompressTaskNoAccess 无权限访问压缩任务
	CompressTaskNoAccess = "CompressTaskNoAccess"
	//DirectoryUsageTaskNotFound 计算目录使用量任务不存在
	DirectoryUsageTaskNotFound = "DirectoryUsageTaskNotFound"
	//DirectoryUsageTaskNoAccess 无权限访问计算目录使用量任务
	DirectoryUsageTaskNoAccess = "DirectoryUsageTaskNoAccess"
	//InvalidDirectoryUsageTaskID 计算目录使用量任务ID不合法
	InvalidDirectoryUsageTaskID = "InvalidArgument.DirectoryUsageTaskID"
	//UnsupportedCompressFileType 不支持的压缩文件类型
	UnsupportedCompressFileType = "UnsupportedCompressFileType"
	//SharedDirectoryExisting 共享目录已存在
	SharedDirectoryExisting = "SharedDirectoryExisting"
	//SharedDirectoryNonexistent 共享目录不存在
	SharedDirectoryNonexistent = "SharedDirectoryNonexistent"
	//SyncTaskNotFound 同步任务不存在
	SyncTaskNotFound = "SyncTaskNotFound"
)

// 压缩器类型
const (
	NONE = "NONE"
	GZIP = "GZIP"
	ZSTD = "ZSTD"
)

// 文件类型
const (
	FILE   = "FILE"
	FOLDER = "FOLDER"
	Batch  = "BATCH"
)

// 操作类型
const (
	UPLOAD     = "UPLOAD"
	DOWNLOAD   = "DOWNLOAD"
	DELETE     = "DELETE"
	MOVE       = "MOVE"
	MKDIR      = "MKDIR"
	COPY       = "COPY"
	COPY_RANGE = "COPY_RANGE"
	COMPRESS   = "COMPRESS"
	CREATE     = "CREATE"
	LINK       = "LINK"
	READ_AT    = "READ_AT"
	WRITE_AT   = "WRITE_AT"
	TRUNCATE   = "TRUNCATE"
)

// 3D云应用
const (
	PlatformLinux          = "LINUX"
	PlatformWindows        = "WINDOWS"
	ScriptRunnerPowershell = "powershell"

	InvalidArgumentSoftwareId                   = "InvalidArgument.SoftwareId"
	InvalidArgumentHardwareId                   = "InvalidArgument.HardwareId"
	InvalidArgumentHardwareSoftwareZoneNotEqual = "InvalidArgument.HardwareSoftwareZoneNotEqual"
	InvalidArgumentSessionId                    = "InvalidArgument.SessionId"
	InvalidArgumentSessionIds                   = "InvalidArgument.SessionIds"
	InvalidArgumentSessionStatus                = "InvalidArgument.SessionStatus"
	InvalidArgumentPostSessionsRequest          = "InvalidArgument.PostSessionsRequest"
	InvalidArgumentRemoteAppName                = "InvalidArgument.RemoteAppName"
	InvalidArgumentRemoteAppId                  = "InvalidArgument.RemoteAppId"
	InvalidArgumentRemoteAppDir                 = "InvalidArgument.Dir"
	InvalidArgumentRemoteAppArgs                = "InvalidArgument.Args"
	InvalidArgumentLogo                         = "InvalidArgument.Logo"
	InvalidArgumentAdminPostRemoteAppsRequest   = "InvalidArgument.AdminPostRemoteAppsRequest"
	InvalidArgumentAdminPutRemoteAppsRequest    = "InvalidArgument.AdminPutRemoteAppsRequest"
	InvalidArgumentAdminPatchRemoteAppsRequest  = "InvalidArgument.AdminPatchRemoteAppsRequest"
	InvalidArgumentAdminCloseSessionRequest     = "InvalidArgument.AdminCloseSessionRequest"
	InvalidArgumentMountPaths                   = "InvalidArgument.MountPaths"
	InvalidArgumentCpu                          = "InvalidArgument.Cpu"
	InvalidArgumentCpuModel                     = "InvalidArgument.CpuModel"
	InvalidArgumentMem                          = "InvalidArgument.Mem"
	InvalidArgumentGpu                          = "InvalidArgument.Gpu"
	InvalidArgumentGpuModel                     = "InvalidArgument.GpuModel"
	InvalidArgumentHardwareName                 = "InvalidArgument.Name"
	InvalidArgumentSoftwareName                 = "InvalidArgument.Name"
	InvalidArgumentSoftwarePlatform             = "InvalidArgument.Platform"
	InvalidArgumentIcon                         = "InvalidArgument.Icon"
	InvalidArgumentSoftwareInitScript           = "InvalidArgument.InitScript"
	InvalidArgumentImageId                      = "InvalidArgument.ImageId"
	InvalidArgumentInstanceType                 = "InvalidArgument.InstanceType"
	InvalidArgumentInstanceFamily               = "InvalidArgument.InstanceFamily"
	InvalidArgumentNetwork                      = "InvalidArgument.Network"
	InvalidArgumentSessionAdminCloseReason      = "InvalidArgument.SessionAdminCloseReason"
	InvalidArgumentRemoteAppLoginUser           = "InvalidArgument.LoginUser"
	InvalidArgumentUserIds                      = "InvalidArgument.UserIds"
	InvalidArgumentScriptRunner                 = "InvalidArgument.ScriptRunner"
	InvalidArgumentScriptContent                = "InvalidArgument.ScriptContent"
	InvalidArgumentMountPoint                   = "InvalidArgument.MountPoint"
	InvalidArgumentShareDirectory               = "InvalidArgument.ShareDirectory"

	HardwareNotFound   = "HardwareNotFound"
	SoftwareNotFound   = "SoftwareNotFound"
	SessionNotFound    = "SessionNotFound"
	RemoteAppNotFound  = "RemoteAppNotFound"
	MountPointNotFound = "MountPointNotFound"

	ForbiddenSessionUserClose    = "Forbidden.SessionUserClose"
	ForbiddenSessionAdminClose   = "Forbidden.SessionAdminClose"
	ForbiddenSessionUserDelete   = "Forbidden.SessionUserDelete"
	ForbiddenSessionUserStart    = "Forbidden.SessionUserStart"
	ForbiddenSessionUserStop     = "Forbidden.SessionUserStop"
	ForbiddenSessionUserRestart  = "Forbidden.SessionUserRestart"
	ForbiddenSessionAdminStart   = "Forbidden.SessionAdminStart"
	ForbiddenSessionAdminStop    = "Forbidden.SessionAdminStop"
	ForbiddenSessionAdminRestart = "Forbidden.SessionAdminRestart"
	ForbiddenSessionRestore      = "Forbidden.SessionRestore"
	ForbiddenSessionNotReady     = "Forbidden.SessionNotReady"

	InternalErrorRunInstanceFailed = "InternalError.RunInstanceFailed"
)

// License相关错误码
const (
	AppIdNotFound      = "AppIdNotFound"
	ManageIdNotFound   = "ManageIdNotFound"
	LicenseIdNotFound  = "LicenseIdNotFound"
	ScIdNotFound       = "ScIdNotFound"
	ModuleIdNotFound   = "ModuleIdNotFound"
	InvalidParams      = "InvalidParams"
	InvalidOs          = "InvalidArgument.Os"
	InvalidComputeRule = "InvalidArgument.ComputeRule"
	InvalidDesc        = "InvalidArgument.Desc"
	InvalidProvider    = "InvalidArgument.Provider"
	InvalidMacAddr     = "InvalidArgument.MacAddr"
	InvalidToolPath    = "InvalidArgument.ToolPath"
	InvalidPort        = "InvalidArgument.Port"
	InvalidLicenseUrl  = "InvalidArgument.LicenseUrl"
	InvalidLicenseNum  = "InvalidArgument.LicenseNum"
	InvalidWeight      = "InvalidArgument.Weight"
	InvalidAuth        = "InvalidArgument.Auth"
	InvalidLicenseType = "InvalidArgument.LicenseType"
	InvalidModuleName  = "InvalidArgument.ModuleName"
	InvalidModuleNum   = "InvalidArgument.ModuleNum"
	InvalidStatus      = "InvalidArgument.Status"
	InvalidPageIndex   = "InvalidArgument.PageIndex"
)

// 商品
const (
	InvalidArgumentMerchandiseName                  = "InvalidArgument.Name"
	InvalidArgumentMerchandiseChargeType            = "InvalidArgument.ChargeType"
	InvalidArgumentUnitPrice                        = "InvalidArgument.UnitPrice"
	InvalidArgumentMerchandiseQuantityUnit          = "InvalidArgument.QuantityUnit"
	InvalidArgumentMerchandiseFormula               = "InvalidArgument.Formula"
	InvalidArgumentMerchandiseYSProduct             = "InvalidArgument.YSProduct"
	InvalidArgumentOutResourceId                    = "InvalidArgument.OutResourceId"
	InvalidArgumentMerchandisePublishState          = "InvalidArgument.PublishState"
	InvalidArgumentMerchandiseId                    = "InvalidArgument.MerchandiseId"
	InvalidArgumentAccountId                        = "InvalidArgument.AccountId"
	InvalidArgumentResourceId                       = "InvalidArgument.ResourceId"
	InvalidArgumentOrderId                          = "InvalidArgument.OrderId"
	InvalidArgumentOrderIdempotentId                = "InvalidArgument.IdempotentId"
	InvalidArgumentOrderQuantity                    = "InvalidArgument.Quantity"
	InvalidArgumentOrderChargeType                  = "InvalidArgument.ChargeType"
	InvalidArgumentOrderStartTime                   = "InvalidArgument.StartTime"
	InvalidArgumentOrderEndTime                     = "InvalidArgument.EndTime"
	InvalidArgumentOrderStartTimeAfterEndTime       = "InvalidArgument.StartTimeAfterEndTime"
	InvalidArgumentOrderStartEndTimeIntervalTooLong = "InvalidArgument.StartEndTimeIntervalTooLong"
	InvalidArgumentOrderIsFirst                     = "InvalidArgument.IsFirst"
	InvalidArgumentOrderIsFinished                  = "InvalidArgument.IsFinished"

	MerchandiseNotFound  = "MerchandiseNotFound"
	OrderNotFound        = "OrderNotFound"
	SpecialPriceNotFound = "SpecialPriceNotFound"

	ForbiddenPatchOrderOnPrePaidOrder = "ForbiddenPatchOrderOnPrePaidOrder"
	ForbiddenPublishMerchandise       = "ForbiddenPublishMerchandise"
	ForbiddenUnpublishMerchandise     = "ForbiddenUnpublishMerchandise"
	ForbiddenSoftwareUser             = "ForbiddenSoftwareUser"
	ForbiddenHardwareUser             = "ForbiddenHardwareUser"

	MerchandiseAlreadyExist   = "MerchandiseAlreadyExist"
	SpecialPriceAlreadyExist  = "SpecialPriceAlreadyExist"
	OrderAlreadyExist         = "OrderAlreadyExist"
	PostPaidOrderAlreadyExist = "PostPaidOrderAlreadyExist"
)

const (
	InvalidArgumentChargeParams          = "InvalidArgument.ChargeParams"
	InternalErrorInvalidAccountId        = "InternalError.InvalidAccountId"
	InvalidAccountStatusNotEnoughBalance = "InvalidAccountStatus.NotEnoughBalance"
	InvalidAccountStatusFrozen           = "InvalidAccountStatus.Frozen"
)

// 代支付
const (
	InvalidArgumentPayBy          = "InvalidArgument.PayBy"
	PayByTokenExpire              = "PayBy.TokenExpire"
	InvalidArgumentPayBySignature = "InvalidArgument.PayBySignature"
)

// cad_converter错误码
const (
	// InvalidArgumentPath 参数错误, path路径不合法，不允许有../之类的内容
	InvalidArgumentPath = "InvalidArgument.Path"
	// InvalidArgumentOverwrite 合法值：true、false
	InvalidArgumentOverwrite = "InvalidArgument.Overwrite"
	// ForbiddenSize 超过最大限制，单文件最大100MB
	ForbiddenSize = "Forbidden.Size"
	// FileExisted 文件已经存在
	FileExisted = "FileExisted"
	// InvalidArgumentSizeTooLarge 文件size太大
	InvalidArgumentSizeTooLarge = "InvalidArgument.SizeTooLarge"
	// InvalidArgumentSize size参数不合法
	InvalidArgumentSize = "InvalidArgument.Size"
	// UploadIdNotFound upload_id不存在
	UploadIdNotFound = "UploadIdNotFound"
	// InvalidArgumentLengthTooLarge 分片长度太大
	InvalidArgumentLengthTooLarge = "InvalidArgument.LengthTooLarge"
	// InvalidArgumentOffset offset超过文件总长度，或者是小于0
	InvalidArgumentOffset = "InvalidArgument.Offset"
	// InvalidArgumentLength length小于0或加上Length后，超过文件总长度
	InvalidArgumentLength = "InvalidArgument.Length"
	// PathNotMachUploadId 文件路径和uploadinit时传的path路径不一致
	PathNotMachUploadId = "PathNotMachUploadId"

	// InputPathNotFound 输入文件路径找不到
	InputPathNotFound = "InputPathNotFound"
	// InvalidArgumentTargetFormat 目标格式不支持
	InvalidArgumentTargetFormat = "InvalidArgument.TargetFormat"
	// InvalidArgumentAutoDelete 自动删除参数不合法,最大604800秒
	InvalidArgumentAutoDelete = "InvalidArgument.AutoDelete"
	// InvalidArgumentPageOffset 分页初始值，默认0，范围 >=0
	InvalidArgumentPageOffset = "InvalidArgument.PageOffset"
	// InvalidArgumentPageSize 分页大小，默认1000，范围1～1000
	InvalidArgumentPageSize = "InvalidArgument.PageSize"
	// InvalidArgumentState 指定的状态不合法
	InvalidArgumentState = "InvalidArgument.State"
	// JobIdNotFound 作业ID未找到
	JobIdNotFound = "JobIdNotFound"
	// OperationConflict 操作冲突
	OperationConflict = "OperationConflict"
	// ForbiddenTotalJobCount 提交作业数量过多
	ForbiddenTotalJobCount = "Forbidden.TotalJobCount"

	// UserIDNotFound 用户ID未找到
	UserIDNotFound = "UserIDNotFound"
	// InvalidArgumentNumber 范围 1～
	InvalidArgumentNumber = "InvalidArgument.Number"
)
