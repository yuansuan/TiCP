package model

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

const TableName = "monitor_node"

// NodeInfo 节点信息
type NodeInfo struct {
	Id              snowflake.ID `json:"id" db:"id" xorm:"pk BIGINT(20)"`
	NodeName        string       `json:"node_name" db:"node_name" xorm:"not null comment('节点名称') VARCHAR(64)"`
	NodeType        string       `json:"node_type" db:"node_type" xorm:"not null comment('节点类型') VARCHAR(16)"`
	SchedulerStatus string       `json:"scheduler_status" db:"scheduler_status" xorm:"not null comment('调度器状态（原始的）') VARCHAR(64)"`
	Status          string       `json:"status" db:"status" xorm:"not null comment('调度器状态（加工后的）') VARCHAR(64)"`
	QueueName       string       `json:"queue_name" db:"queue_name" xorm:"not null comment('所属队列') VARCHAR(16)"`
	PlatformName    string       `json:"platform_name" db:"platform_name" xorm:"not null default '' comment('标识') VARCHAR(64)"`
	TotalCoreNum    int          `json:"total_core_num" db:"total_core_num" xorm:"not null default 0 comment('总核数') INT(11)"`
	UsedCoreNum     int          `json:"used_core_num" db:"used_core_num" xorm:"not null default 0 comment('已使用核数') INT(11)"`
	FreeCoreNum     int          `json:"free_core_num" db:"free_core_num" xorm:"not null default 0 comment('剩余的核数') INT(11)"`
	TotalMem        int          `json:"total_mem" db:"total_mem" xorm:"not null default 0 comment('总内存空间') INT(11)"`
	UsedMem         int          `json:"used_mem" db:"used_mem" xorm:"not null default 0 comment('已使用的内存空间') INT(11)"`
	FreeMem         int          `json:"free_mem" db:"free_mem" xorm:"not null default 0 comment('剩余内存空间') INT(11)"`
	AvailableMem    int          `json:"available_mem" db:"available_mem" xorm:"not null default 0 comment('可以使用的内存空间') INT(11)"`
	CreateTime      time.Time    `json:"create_time" db:"create_time" xorm:"not null comment('创建时间') DATETIME updated"`
	UpdateTime      time.Time    `json:"update_time" db:"update_time" xorm:"not null comment('更新时间时间') DATETIME updated"`
}

func (n *NodeInfo) TableName() string {
	return TableName
}
