package company

// CompanyUserConfigKeys define company-user config keys
var CompanyUserConfigKeys = struct {
	KeySchedulerCloudCostLimitDay   string
	KeySchedulerCloudCostLimitMonth string
}{
	KeySchedulerCloudCostLimitDay:   "scheduler.cloud.cost_limit.day",   // yuan per day
	KeySchedulerCloudCostLimitMonth: "scheduler.cloud.cost_limit.month", // yuan per month
}

// CompanyConfigKeys define company config keys
var CompanyConfigKeys = struct {
	KeySchedulerTimeBeforeBoost   string
	KeyDownloadSpeedLimit         string
	KeyUploadSpeedLimit           string
	KeyFrontendCloudType          string
	KeyFrontendLiveChatID         string
	KeyAlarmValueOfAccountBalance string
	KeyVisualServiceOn            string
	KeyVisualMaxTerminal          string
}{
	KeySchedulerTimeBeforeBoost:   "scheduler.local.time_before_boost", // seconds
	KeyDownloadSpeedLimit:         "file_sync.limit.download",          // bytes per seconds
	KeyUploadSpeedLimit:           "file_sync.limit.upload",            // bytes per seconds
	KeyFrontendCloudType:          "frontend.cloud_type",               // "public", "mixed"
	KeyFrontendLiveChatID:         "frontend.live_chat_id",             // "8b8738ec"
	KeyAlarmValueOfAccountBalance: "alarm_value.account_balance",       // alarm value of account balance
	KeyVisualServiceOn:            "is_visual_service_on",
	KeyVisualMaxTerminal:          "max_visual_termianl",
}
