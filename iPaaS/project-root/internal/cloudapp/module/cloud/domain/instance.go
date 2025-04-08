package domain

import (
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
)

// Instance 表示一个在第三方云服务启动的实例
type Instance interface {
	// ID 实例的ID
	ID() string

	// Raw 返回实例的原始数据
	Raw() string

	// Status 实例当前状态
	Status() models.InstanceStatus
}
