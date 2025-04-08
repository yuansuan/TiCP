package model

import (
	"time"
)

type UploadInfo struct {
	Id         string    `json:"id" xorm:"pk comment('上传 ID') VARCHAR(255)"`
	UserId     string    `json:"user_id" xorm:"not null comment('请求者用户 ID') VARCHAR(32)"`
	TmpPath    string    `json:"tmp_path" xorm:"not null default '' comment('临时文件路径') MEDIUMTEXT"`
	Path       string    `json:"path" xorm:"not null default '' comment('文件路径') MEDIUMTEXT"`
	Size       int64     `json:"size" xorm:"not null default 0 comment('文件大小') BIGINT(20)"`
	Overwrite  bool      `json:"overwrite" xorm:"not null default 0 comment('是否覆盖') TINYINT(1)"`
	CreateTime time.Time `json:"create_time" xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime time.Time `json:"update_time" xorm:"not null default CURRENT_TIMESTAMP comment('更新时间') DATETIME"`
}
