package dto

type JobTimeLine struct {
	EventName string `json:"name"`     // 名称
	EventTime string `json:"time"`     // 时间
	Progress  int    `json:"progress"` // 进度
}
