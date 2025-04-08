package model

import "time"

type CompressInfo struct {
	Id         string             `json:"id" xorm:"pk comment('压缩ID') VARCHAR(255)"`
	UserId     string             `json:"user_id" xorm:"not null comment('用户ID') VARCHAR(32)"`
	TmpPath    string             `json:"tmp_path" xorm:"not null default '' comment('临时文件路径') MEDIUMTEXT"`
	Paths      string             `json:"paths" xorm:"not null default '' comment('源文件/文件夹路径') MEDIUMTEXT"`
	TargetPath string             `json:"target_path" xorm:"not null default '' comment('目标文件路径') MEDIUMTEXT"`
	BasePath   string             `json:"base_path" xorm:"default '' comment('压缩包起始目录路径') MEDIUMTEXT"`
	ErrorMsg   string             `json:"error_msg" xorm:"default '' comment('错误信息') TEXT"`
	Status     CompressTaskStatus `json:"status" xorm:"not null default 0 comment('状态 -1: 失败 0: 压缩中 1: 成功 2: 已取消') TINYINT(1)"`
	CreateTime time.Time          `json:"create_time" xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime time.Time          `json:"update_time" xorm:"not null default CURRENT_TIMESTAMP comment('更新时间') DATETIME"`
}

type CompressTaskStatus int

const (
	CompressTaskFailed   CompressTaskStatus = -1
	CompressTaskRunning  CompressTaskStatus = 0
	CompressTaskFinished CompressTaskStatus = 1
	CompressTaskCanceled CompressTaskStatus = 2
)
