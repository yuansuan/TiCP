package consts

const (
	Queue            = "QUEUE"
	CpuNum           = "CPU_NUM"
	OriginCpuNum     = "ORIGIN_CPU_NUM"
	MemNum           = "MEM_NUM"
	NodeNum          = "NODE_NUM"
	JobName          = "JOB_NAME"
	DateTime         = "DATE_TIME"
	MainFile         = "YS_MAIN_FILE"
	Platform         = "PLATFORM"
	ForbidDownload   = "FORBID_DOWNLOAD"
	NodeSelector     = "NODE_SELECTOR"
	TextType         = "text"
	ListType         = "list"
	MultipleType     = "multiple"
	NodeSelectorType = "node_selector"

	Semicolon                  = ";"
	DefaultSize                = 30
	WorkDirSuffixDefault       = "default"
	WorkDirDefaultTypeMainFile = "main_file"
	AppStatePublished          = "published"

	TimelineLocalPrefix    = "Local"
	TimelineCloudPrefix    = "Cloud"
	JobAttrKeySubmitEnvs   = "SubmitEnvs"
	JobAttrKeySubmitParams = "SubmitParams"
)

const (
	APIJobStatePending            = "Pending"
	APIJobStateRunning            = "Running"
	APIJobStateTransmitting       = "Transmitting"
	APIJobStateTerminating        = "Terminating"
	APIJobStateTerminated         = "Terminated"
	APIJobStateSuspending         = "Suspending"
	APIJobStateSuspended          = "Suspended"
	APIJobStateCompleted          = "Completed"
	APIJobStateFailed             = "Failed"
	APIJobStateInitiallySuspended = "InitiallySuspended"
)

const (
	JobStateSubmitted   = "Submitted"
	JobStatePending     = "Pending"
	JobStateRunning     = "Running"
	JobStateTerminated  = "Terminated"
	JobStateSuspended   = "Suspended"
	JobStateCompleted   = "Completed"
	JobStateFailed      = "Failed"
	JobStateBursting    = "Bursting"
	JobStateBurstFailed = "BurstFailed"
)

const (
	JobDataStateUploading      = "Uploading"
	JobDataStateUploaded       = "Uploaded"
	JobDataStateUploadFailed   = "UploadFailed"
	JobDataStateDownloading    = "Downloading"
	JobDataStateDownloaded     = "Downloaded"
	JobDataStateDownloadFailed = "DownloadFailed"
)

const (
	JobDeliveryCountUser = "User"
	JobDeliveryCountJob  = "Job"
)

const (
	JobWaitTimeStatisticAvg   = "JobWaitTimeStatisticAvg"
	JobWaitTimeStatisticMax   = "JobWaitTimeStatisticMax"
	JobWaitTimeStatisticTotal = "JobWaitTimeStatisticTotal"
	JobWaitNumStatistic       = "JobWaitNumStatistic"
)

const (
	JobStatisticsQueryTypeApp     = "app"
	JobStatisticsQueryTypeUser    = "user"
	JobStatisticsShowTypeOverview = "overview"
	JobStatisticsShowTypeDetail   = "detail"
)

const (
	JobVisAnalysisResidual = "vis_analysis_residual"
	JobVisAnalysisSnapshot = "vis_analysis_snapshot"
)

const (
	JobQueryUser = "用户"
	JobQueryApp  = "应用"
	JobOverview  = "统计总览"
	JobDetail    = "统计明细"
)
