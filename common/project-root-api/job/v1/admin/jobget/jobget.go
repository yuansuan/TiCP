package jobget

import (
	get "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobget"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
// swagger:model AdminJobGetRequest
type Request struct {
	get.Request `json:",inline"`
}

// Response 返回
// swagger:model AdminJobGetResponse
type Response struct {
	schema.Response `json:",inline"`

	Data *Data `json:"Data,omitempty"`
}

// Data 数据
// swagger:model AdminJobGetData
type Data struct {
	schema.AdminJobInfo `json:",inline"`
}
