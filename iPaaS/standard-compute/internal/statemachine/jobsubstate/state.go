package jobsubstate

type SubState string

func (s SubState) String() string {
	return string(s)
}

const (
	SubStateNull = SubState("")

	PreparingPulling     = SubState("preparing-pulling-image")    // 拉取镜像
	PreparingDownloading = SubState("preparing-downloading-file") // 下载输入文件
	PreparingSubmitting  = SubState("preparing-submitting-job")   // 提交作业到调度器

	PendingWaitingSchedule = SubState("pending-waiting-schedule") // 等待调度器开始执行

	RunningWaitingResult = SubState("running-waiting-result") // 等待运行完成

	CompletingUploading = SubState("completing-uploading-file") // 数据回传

	CompletedAllDone = SubState("completed-all-done") // 作业执行完成

	CanceledByUser    = SubState("canceled-by-user")    // 用户执行的取消作业
	CanceledByTimeout = SubState("canceled-by-timeout") // 由于超时导致的关闭
)
