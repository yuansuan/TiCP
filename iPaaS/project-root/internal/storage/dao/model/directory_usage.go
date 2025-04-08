package model

import "time"

type DirectoryUsage struct {
	Id         string                   `json:"id" xorm:"pk comment('目录用量计算任务id') VARCHAR(255)"`
	UserID     string                   `json:"user_id" xorm:"user_id not null comment('用户id') VARCHAR(255)"`
	Path       string                   `json:"path" xorm:"not null default '' comment('目录路径') MEDIUMTEXT"`
	Size       int64                    `json:"size" xorm:"not null default 0 comment('目录大小 单位为字节') BIGINT(20)"`
	LogicSize  int64                    `json:"logic_size" xorm:"not null default 0 comment('目录大小 单位为字节 软链接本身大小') BIGINT(20)"`
	Status     DirectoryUsageTaskStatus `json:"status" xorm:"not null default 0 comment('状态 -1: 失败 0: 计算中 1: 成功 2: 已取消') TINYINT(1)"`
	ErrMsg     string                   `json:"err_msg" xorm:"not null default '' comment('错误信息') TEXT"`
	CreateTime time.Time                `json:"create_time" xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime time.Time                `json:"update_time" xorm:"not null default CURRENT_TIMESTAMP comment('更新时间') DATETIME"`
}

type DirectoryUsageTaskStatus int

const (
	DirectoryUsageTaskFailed      DirectoryUsageTaskStatus = -1
	DirectoryUsageTaskCalculating DirectoryUsageTaskStatus = 0
	DirectoryUsageTaskFinished    DirectoryUsageTaskStatus = 1
	DirectoryUsageTaskCanceled    DirectoryUsageTaskStatus = 2
)
