package jobcpuusage

import (
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcpuusage"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
type Request struct {
	jobcpuusage.Request `json:",inline"`
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
