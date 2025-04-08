package model

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type App struct {
	ID                snowflake.ID      `xorm:"id BIGINT(20) pk"`
	OutAppID          string            `xorm:"out_app_id VARCHAR(255)"`
	OutID             string            `xorm:"out_id VARCHAR(255)"`
	OutName           string            `xorm:"out_name VARCHAR(255)"`
	Name              string            `xorm:"name VARCHAR(255)"`
	VersionNum        string            `xorm:"version_num VARCHAR(255)"`
	Type              string            `xorm:"type VARCHAR(255)"`
	ComputeType       string            `xorm:"compute_type VARCHAR(64)"`
	QueueNames        []string          `xorm:"queue_names Text"`
	State             string            `xorm:"state VARCHAR(64)"`
	Image             string            `xorm:"image VARCHAR(64)"`
	BinPath           map[string]string `xorm:"bin_path TEXT"`
	SchedulerParam    map[string]string `xorm:"scheduler_param TEXT"`
	LicenseManagerId  string            `xorm:"license_manager_id VARCHAR(64)"`
	EnableResidual    bool              `xorm:"enable_residual TINYINT(4)"`
	ResidualLogParser string            `xorm:"residual_log_parser VARCHAR(64)"`
	EnableSnapshot    bool              `xorm:"enable_snapshot TINYINT(4)"`
	Script            string            `xorm:"script MEDIUMTEXT"`
	Content           string            `xorm:"content MEDIUMTEXT"`
	Icon              string            `xorm:"icon MEDIUMTEXT"`
	Description       string            `xorm:"description TEXT"`
	DocType           string            `xorm:"doc_type VARCHAR(64)"`
	DocContent        string            `xorm:"doc_content TEXT"`
	CreateTime        time.Time         `xorm:"create_time DATETIME created"`
	UpdateTime        time.Time         `xorm:"update_time DATETIME updated"`
}

const AppTableName = "app_template"

func (App) TableName() string {
	return AppTableName
}
