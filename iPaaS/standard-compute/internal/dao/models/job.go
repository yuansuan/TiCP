package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/filesyncstate"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/jobstate"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/statemachine/jobsubstate"
)

type Job struct {
	Id           int64
	IdempotentId string

	State        jobstate.State
	SubState     jobsubstate.SubState
	IsUserFailed bool

	FileSyncState filesyncstate.State

	DownloadCurrentSize int64
	DownloadTotalSize   int64
	UploadCurrentSize   int64
	UploadTotalSize     int64

	StateReason string

	//请求的
	Queue         string
	Priority      int64
	RequestCores  int64
	CoresPerNode  int64 // 含义；单机核数
	RequestMemory int64
	AllocType     string // average; manual; 其他

	// "local" "image"
	AppMode          AppMode
	AppPath          string
	SingularityImage string

	Inputs               string
	Output               string
	EnvVars              string
	Command              string
	SchedulerSubmitFlags string

	Workspace string
	Script    string
	Stdout    string
	Stderr    string

	// 原地计算
	IsOverride bool
	WorkDir    string

	//实际的
	OriginJobId       string
	AllocCores        int64
	AllocMemory       int64
	OriginState       string
	ExitCode          string
	ExecutionDuration int64
	PendingTime       *time.Time
	RunningTime       *time.Time
	CompletingTime    *time.Time
	CompletedTime     *time.Time

	ExecHosts   string
	ExecHostNum int64

	ControlBitTerminate bool
	IsTimeout           bool
	Timeout             int64

	Webhook string

	CustomStateRule string

	CreateTime time.Time `xorm:"created"`
	UpdateTime time.Time `xorm:"updated"`
}

func (j *Job) TableName() string {
	return "sc_job"
}

type AppMode string

const (
	ImageAppMode AppMode = "image"
	LocalAppMode AppMode = "local"
)

var appModeMap = map[string]AppMode{
	string(ImageAppMode): ImageAppMode,
	string(LocalAppMode): LocalAppMode,
}

func (m AppMode) String() string {
	return string(m)
}

func StrToAppMode(s string) (AppMode, error) {
	mode, exist := appModeMap[s]
	if !exist {
		return "", fmt.Errorf("app mode [%s] not found", s)
	}

	return mode, nil
}

func (j *Job) JobID() string {
	return strconv.FormatInt(j.Id, 10)
}

func (j *Job) ToHTTPModel() *v20230530.JobInHPC {
	m := &v20230530.JobInHPC{
		ID:          snowflake.ParseInt64(j.Id).String(),
		Environment: mustRestoreEnv(j.EnvVars),
		Command:     j.Command,
		Override: v20230530.JobInHPCOverride{
			Enable:  j.IsOverride,
			WorkDir: j.WorkDir,
		},
		Queue: j.Queue,
		Resource: v20230530.JobInHPCResource{
			Cores:     int(j.RequestCores),
			AllocType: j.AllocType,
		},
		Inputs:            mustRestoreInputs(j.Inputs),
		Output:            mustRestoreOutput(j.Output),
		CustomStateRule:   mustRestoreCustomStateRule(j.CustomStateRule),
		SchedulerID:       j.OriginJobId,
		Status:            j.State,
		IsUserFailed:      j.IsUserFailed,
		StateReason:       j.StateReason,
		PendingTime:       j.PendingTime,
		RunningTime:       j.RunningTime,
		CompletingTime:    j.CompletingTime,
		CompletedTime:     j.CompletedTime,
		AllocCores:        int(j.AllocCores),
		ExitCode:          j.ExitCode,
		ExecutionDuration: int(j.ExecutionDuration),
		DownloadProgress: v20230530.JobInHPCProgress{
			Total:   int(j.DownloadTotalSize),
			Current: int(j.DownloadCurrentSize),
		},
		Priority:     int(j.Priority),
		ExecHosts:    j.ExecHosts,
		ExecHostsNum: int(j.ExecHostNum),
	}
	if j.AppMode == ImageAppMode {
		m.Application = fmt.Sprintf("image:%s", j.SingularityImage)
	} else if j.AppMode == LocalAppMode {
		m.Application = fmt.Sprintf("local:%s", j.AppPath)
	}

	return m
}

func mustRestoreEnv(envsStr string) map[string]string {
	envs := make([]string, 0)
	res := make(map[string]string)

	var err error
	if err = jsoniter.UnmarshalFromString(envsStr, &envs); err != nil {
		log.Warnf("unmarshal from string %s to []string failed, %v", envsStr, err)
		return res
	}

	for _, env := range envs {
		fields := strings.Split(env, "=")

		if len(fields) != 2 {
			continue
		}

		res[fields[0]] = fields[1]
	}

	return res
}

func mustRestoreInputs(inputsStr string) []v20230530.JobInHPCInputStorage {
	inputs := make([]v20230530.JobInHPCInputStorage, 0)

	if err := jsoniter.UnmarshalFromString(inputsStr, &inputs); err != nil {
		log.Warnf("unmarshal from string %s to []model.InputStorage failed, %v", inputsStr, err)
		return inputs
	}

	return inputs
}

func mustRestoreOutput(outputStr string) *v20230530.JobInHPCOutputStorage {
	output := new(v20230530.JobInHPCOutputStorage)

	if err := jsoniter.UnmarshalFromString(outputStr, output); err != nil {
		log.Warnf("unmarshal from string %s to *model.OutputStorage failed, %v", outputStr, err)
		return output
	}

	return output
}

func mustRestoreCustomStateRule(customStateRuleStr string) *v20230530.JobInHPCCustomStateRule {
	customStateRule := new(v20230530.JobInHPCCustomStateRule)

	if err := jsoniter.UnmarshalFromString(customStateRuleStr, customStateRule); err != nil {
		log.Warnf("unmarshal from string %s to *model.CustomStateRule failed, %v", customStateRuleStr, err)
		return nil
	}

	return customStateRule
}
