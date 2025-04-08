package models

import (
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

// Application 应用
type Application struct {
	ID                 snowflake.ID `json:"id" xorm:"id not null pk BIGINT(20)"`
	Name               string       `xorm:"varchar(255) not null comment('display of the application, such as: Abaqus 6.1.5')" json:"name"`
	Type               string       `xorm:"varchar(255) not null comment('real name of the application, such as: Abaqus, used to classify applications without version')" json:"type"`
	Version            string       `xorm:"varchar(32) not null comment('version of the application, such as: 6.1.5')" json:"version"`
	AppParamsVersion   int          `xorm:"int(11) not null comment('app_params version')" json:"app_params_version"`
	Image              string       `xorm:"varchar(128) not null comment('image名称')" json:"image"`
	Endpoint           string       `xorm:"varchar(255) not null comment('超算中心endpoint')" json:"endpoint"`
	Command            string       `xorm:"text not null comment('提交命令')" json:"command"`
	PublishStatus      string       `xorm:"varchar(32) not null default 'unpublished' comment('发布状态')" json:"publish_status"`
	Description        string       `xorm:"text comment('应用描述')" json:"description"`
	IconUrl            string       `xorm:"varchar(128) not null comment('应用图标')" json:"icon_url"`
	CoresMaxLimit      int64        `xorm:"bigint(20) not null default '0' comment('cores_max_limit')" json:"cores_max_limit"`
	CoresPlaceholder   string       `xorm:"varchar(256) not null default '' comment('cores_placeholder')" json:"cores_placeholder"`
	FileFilterRule     string       `xorm:"varchar(255) default '{"result": "\\\\.dat$","model": "\\\\.(jou|cas)$","log": "\\\\.(sta|dat|msg|out|log)$","middle": "\\\\.(com|prt)$"}' comment('文件过滤规则')" json:"file_filter_rule"`
	ResidualEnable     bool         `xorm:"tinyint(1) default '0' comment('残差图是否开启')" json:"residual_enable"`
	ResidualLogRegexp  string       `xorm:"varchar(255) default 'stdout.log' comment('残差图文件')" json:"residual_log_regexp"`
	ResidualLogParser  string       `xorm:"varchar(255) default '' comment('残差图解析器')" json:"residual_log_parser"`
	MonitorChartEnable bool         `xorm:"tinyint(1) default '0' comment('监控图表是否开启')" json:"monitor_chart_enable"`
	MonitorChartRegexp string       `xorm:"varchar(255) default '.*\\.out' comment('监控图表文件规则')" json:"monitor_chart_regexp"`
	MonitorChartParser string       `xorm:"varchar(255) default '' comment('监控图表解析器')" json:"monitor_chart_parser"`
	SnapshotEnable     bool         `xorm:"tinyint(1) default '0' comment('云图是否开启')" json:"snapshot_enable"`
	BinPath            string       `xorm:"text comment('bin_path')" json:"bin_path"`
	ExtentionParams    string       `xorm:"text comment('扩展参数')" json:"extention_params"`
	LicManagerId       snowflake.ID `xorm:"BIGINT(20) default 0 comment('LicenseManager id, 如果为空代表免费软件')" json:"lic_manager_id"`
	NeedLimitCore      bool         `xorm:"tinyint(1) default '0' comment('是否需要限制核数')" json:"need_limit_core"`
	SpecifyQueue       string       `xorm:"varchar(255) default '' comment('指定队列')" json:"specify_queue"`
	CreateTime         time.Time    `xorm:"datetime not null default comment('创建时间') created" json:"create_time"`
	UpdateTime         time.Time    `xorm:"datetime not null default comment('修改时间') updated" json:"update_time"`
}
