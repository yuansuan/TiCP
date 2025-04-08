package models

import (
	"encoding/json"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

type PreSchedule struct {
	// 作业基本信息
	ID snowflake.ID `json:"id" xorm:"pk id comment('预调度ID') BIGINT(20)"`

	// 参数信息
	Params          string `json:"params" xorm:"comment('用户参数') TEXT"`
	ExpectedMinCpus int64  `json:"expected_min_cpus" xorm:"not null default 1 comment('期望的最小核数') BIGINT(20)"`
	ExpectedMaxCpus int64  `json:"expected_max_cpus" xorm:"not null default 1 comment('期望的最大核数') BIGINT(20)"`
	ExpectedMemory  int64  `json:"expected_memory" xorm:"not null default 1 comment('期望的内存数') BIGINT(20)"`
	Shared          bool   `json:"shared" xorm:"not null default 0 comment('是否共享节点') TINYINT(1)"`
	Fixed           bool   `json:"fixed" xorm:"not null default 0 comment('是否固定分区') TINYINT(1)"`

	// 作业运行信息
	Zone    string `json:"zone" xorm:"not null default '' comment('预调度的分区') VARCHAR(64)"`
	Command string `json:"command" xorm:"not null comment('作业实际执行命令') TEXT"`
	WorkDir string `json:"work_dir" xorm:"not null default '' comment('工作目录') VARCHAR(255)"`

	// 应用信息
	AppID   snowflake.ID `json:"app_id" xorm:"app_id not null comment('计算应用ID') BIGINT(20)"`
	AppName string       `json:"app_name" xorm:"not null comment('计算应用名') VARCHAR(255)"`
	Envs    string       `json:"envs" xorm:"comment('环境变量') TEXT"`

	// 是否已使用
	Used bool `json:"used" xorm:"not null default 0 comment('是否已使用') TINYINT(1)"`

	// 时间信息
	CreateTime time.Time `json:"create_time" xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME created"`
	UpdateTime time.Time `json:"update_time" xorm:"not null default CURRENT_TIMESTAMP comment('更新时间') DATETIME updated"`
}

func (m *PreSchedule) GetResource() (*int, *int, *int) {
	// int64 -> *int
	mincpus := int(m.ExpectedMinCpus)
	maxcpus := int(m.ExpectedMaxCpus)
	memory := int(m.ExpectedMemory)
	return &mincpus, &maxcpus, &memory
}

func (m *PreSchedule) GetEnvVars() (map[string]string, error) {
	if m.Envs == "" {
		return nil, nil
	}
	envs := make(map[string]string)
	if err := json.Unmarshal([]byte(m.Envs), &envs); err != nil {
		return nil, err
	}
	return envs, nil
}
