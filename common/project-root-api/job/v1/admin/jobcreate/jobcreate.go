package jobcreate

import (
	job "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcreate"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
// swagger:model AdminJobCreateRequest
type Request struct {
	job.Request             `json:",inline"`
	Queue                   string            `json:"Queue"`                             // 指定作业运行的队列
	JobSchedulerSubmitFlags map[string]string `json:"JobSchedulerSubmitFlags,omitempty"` //自定义调度器提交参数  {"-o": "stdout.log", "-e", "stderr.log", "-l": "select=xxx"}
}

// Response 返回
// swagger:model AdminJobCreateResponse
type Response struct {
	schema.Response `json:",inline"`
	Data            *Data `json:"Data,omitempty"`
}

// Data 数据
// swagger:model AdminJobCreateData
type Data struct {
	JobID string `json:"JobID"`
}
