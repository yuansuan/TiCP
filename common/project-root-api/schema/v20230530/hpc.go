package v20230530

import (
	"time"

	"github.com/yuansuan/ticp/common/project-root-api/hpc/filesyncstate"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/jobstate"
)

type JobInHPC struct {
	ID                string                   `json:"ID,omitempty"`
	Application       string                   `json:"Application,omitempty"`
	Environment       map[string]string        `json:"Environment,omitempty"`
	Command           string                   `json:"Command,omitempty"`
	Override          JobInHPCOverride         `json:"Override,omitempty"`
	Queue             string                   `json:"Queue,omitempty"`
	Resource          JobInHPCResource         `json:"Resource,omitempty"`
	Inputs            []JobInHPCInputStorage   `json:"Inputs,omitempty"`
	Output            *JobInHPCOutputStorage   `json:"Output,omitempty"`
	CustomStateRule   *JobInHPCCustomStateRule `json:"CustomStateRule,omitempty"`
	SchedulerID       string                   `json:"SchedulerID,omitempty"`
	Status            jobstate.State           `json:"Status"`
	IsUserFailed      bool                     `json:"IsUserFailed,omitempty"`
	FileSyncState     filesyncstate.State      `json:"FileSyncState,omitempty"`
	StateReason       string                   `json:"StateReason,omitempty"`
	PendingTime       *time.Time               `json:"PendingTime,omitempty"`
	RunningTime       *time.Time               `json:"RunningTime,omitempty"`
	CompletingTime    *time.Time               `json:"CompletingTime,omitempty"`
	CompletedTime     *time.Time               `json:"CompletedTime,omitempty"`
	AllocCores        int                      `json:"AllocCores,omitempty"`
	ExitCode          string                   `json:"ExitCode,omitempty"`
	ExecutionDuration int                      `json:"ExecutionDuration,omitempty"`
	DownloadProgress  JobInHPCProgress         `json:"DownloadProgress,omitempty"`
	UploadProgress    JobInHPCProgress         `json:"UploadProgress,omitempty"`
	Priority          int                      `json:"Priority,omitempty"`
	ExecHosts         string                   `json:"ExecHosts,omitempty"`
	ExecHostsNum      int                      `json:"ExecHostsNum,omitempty"`
}

type JobInHPCProgress struct {
	Total   int `json:"Total,omitempty"`
	Current int `json:"Current,omitempty"`
}

// JobInHPCOverride 原地计算相关，启用原地计算时，Enable设置为true，WorkDir遵循url格式，例: "schema://host/path"
type JobInHPCOverride struct {
	Enable  bool   `json:"Enable,omitempty"`
	WorkDir string `json:"WorkDir,omitempty"`
}

type JobInHPCResource struct {
	Cores        int    `json:"Cores,omitempty"`        //申请的核数
	CoresPerNode *int   `json:"CoresPerNode,omitempty"` //自定义单机核数
	AllocType    string `json:"AllocType"`              //CPU资源的分配方式
}

// JobInHPCInputStorage example
// "src": "http://192.168.56.107:10000/testA",
// "dst": "/testA",
// "type": "hpc_storage"
type JobInHPCInputStorage struct {
	Src  string      `json:"Src,omitempty"`
	Dst  string      `json:"Dst,omitempty"`
	Type StorageType `json:"Type,omitempty"`
}

type StorageType string

const (
	HPCStorageType   StorageType = "hpc_storage"
	CloudStorageType StorageType = "cloud_storage"
)

// JobInHPCOutputStorage example
// "dst": "http://192.168.56.107:10000/dirA",
// "type": "hpc_storage",
type JobInHPCOutputStorage struct {
	Dst           string      `json:"Dst,omitempty"`
	Type          StorageType `json:"Type,omitempty"`
	NoNeededPaths string      `json:"NoNeededPaths,omitempty"`
	NeededPaths   string      `json:"NeededPaths,omitempty"`
}

type JobInHPCCustomStateRule struct {
	KeyStatement string `json:"KeyStatement,omitempty"`
	// ResultState 包含KeyStatement关键字时，认为的结果状态，仅有两种可选择[ completed | failed ]
	ResultState string `json:"ResultState,omitempty"`
}

type Resource struct {
	Cpu           int64   `json:"Cpu"` //空闲cpu
	TotalCpu      int64   `json:"TotalCpu"`
	CoresPerNode  int64   `json:"CoresPerNode"`
	IdleNodeNum   int64   `json:"IdleNodeNum"`
	AllocNodeNum  int64   `json:"AllocNodeNum"` //已使用的节点
	TotalNodeNum  int64   `json:"TotalNodeNum"`
	Memory        int64   `json:"Memory"`      //空闲内存 单位mb
	TotalMemory   int64   `json:"TotalMemory"` //单位mb
	IsDefault     bool    `json:"IsDefault"`
	ReservedCores int64   `json:"ReservedCores"` // 预留核数
	StorageUsage  float64 `json:"StorageUsage"`  //使用存储 单位GB
	StorageLimit  float64 `json:"StorageLimit"`
}

type CpuUsage struct {
	JobID           string             `json:"JobID"`
	AverageCpuUsage float64            `json:"AverageCpuUsage"`
	NodeUsages      map[string]float64 `json:"NodeUsages"`
}
