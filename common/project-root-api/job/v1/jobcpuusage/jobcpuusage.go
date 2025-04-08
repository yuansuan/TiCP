package jobcpuusage

import (
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
type Request struct {
	JobID string `uri:"JobID" binding:"required"` //作业ID
}

// Response 返回
type Response struct {
	schema.Response `json:",inline"`

	Data *Data `json:"Data,omitempty"`
}

// Data 数据
type Data struct {
	schema.JobCpuUsage `json:",inline"`
}
