package jobpreschedule

import (
	job "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcreate"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
type Request struct {
	Params Params   `json:"Params" binding:"required"` //预调度作业参数
	Zones  []string `json:"Zones"`                     //预选分区列表
	Shared bool     `json:"Shared"`                    //是否共享节点(同提交作业的 NoRound 字段)，如果为true，则单节点核数不进行取整，节点剩余核数可供其他作业使用，字段仅限内部用户使用
	Fixed  bool     `json:"Fixed"`                     //是否固定分区，为true时，只能在指定的分区中调度
}

type Params struct {
	Application job.Application   `json:"Application" binding:"required"` //求解软件信息
	Resource    *Resource         `json:"Resource" binding:"required"`    //作业需要的期望资源
	EnvVars     map[string]string `json:"EnvVars"`                        //环境变量
}

type Resource struct {
	MinCores *int `json:"MinCores" binding:"required"` //期望的最小核数
	MaxCores *int `json:"MaxCores" binding:"required"` //期望的最大核数
	Memory   *int `json:"Memory"`                      //期望的内存数，单位为M，暂未起作用
}

// Response 响应
type Response struct {
	schema.Response `json:",inline"`
	Data            *Data `json:"Data,omitempty"`
}

// Data 数据
type Data struct {
	ScheduleID string `json:"ScheduleID"`
	Workdir    string `json:"Workdir"`
}
