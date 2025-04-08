package queue

type Queue struct {
	// 任务排队数量
	JobPendingNum int64
	// 任务运行数量
	JobRunningNum int64
}
