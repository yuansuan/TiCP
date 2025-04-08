package model

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

// Job 作业信息
type Job struct {
	Id            snowflake.ID
	AppId         snowflake.ID
	UserId        snowflake.ID
	JobSetId      snowflake.ID
	ProjectId     snowflake.ID
	OutJobId      string
	RealJobId     string
	UploadTaskId  string
	Type          string
	Name          string
	Queue         string
	State         string
	RawState      string
	DataState     string
	ExitCode      string
	AppName       string
	UserName      string
	JobSetName    string
	ProjectName   string
	ClusterName   string
	BurstNum      int
	VisAnalysis   map[string]bool
	Priority      int
	CpusAlloc     int
	MemAlloc      int
	ExecDuration  int
	ExecHostNum   int
	Reason        string
	WorkDir       string
	ExecHosts     string
	SubmitTime    time.Time
	PendTime      time.Time
	StartTime     time.Time
	EndTime       time.Time
	TerminateTime time.Time
	SuspendTime   time.Time
	CreateTime    time.Time `xorm:"create_time DATETIME created"`
	UpdateTime    time.Time `xorm:"update_time DATETIME updated"`
}

func (j *Job) TableName() string {
	return "job"
}

type StatisticsJob struct {
	*Job    `xorm:"extends"`
	CPUTime float64 `xorm:"cpu_time"`
}
