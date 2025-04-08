package model

import "time"

// SharedDirectory 共享目录数据
type SharedDirectory struct {
	ID             int64     `json:"id" xorm:"id pk autoincr comment('主键id') BIGINT(20)"`
	UserID         string    `json:"user_id" xorm:"user_id not null comment('用户id') VARCHAR(255)"`
	Path           string    `json:"path" xorm:"path not null comment('指定路径') MEDIUMTEXT"`
	SharedUserName string    `json:"shared_user_name" xorm:"shared_user_name not null comment('用户名') VARCHAR(255)"`
	SharedPassword string    `json:"shared_password" xorm:"shared_password not null comment('密码') VARCHAR(255)"`
	SharedHost     string    `json:"shared_host" xorm:"shared_host not null comment('共享主机地址') VARCHAR(255)"`
	SharedSrc      string    `json:"shared_src" xorm:"shared_src not null comment('共享目录路径') VARCHAR(255)"`
	IsDeleted      int       `json:"is_deleted" xorm:"is_deleted not null default 0 comment('是否删除') TINYINT(1)"`
	CreateTime     time.Time `json:"create_time" xorm:"create_time not null created default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	UpdateTime     time.Time `json:"update_time" xorm:"update_time not null updated default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}
