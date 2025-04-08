package joblist

import (
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
// swagger:model JobListRequest
type Request struct {
	JobState   string `form:"JobState"`
	Zone       string `form:"Zone"`
	PageOffset *int64 `form:"PageOffset"`
	PageSize   *int64 `form:"PageSize"`
}

// Response 返回
// swagger:model JobListResponse
type Response struct {
	schema.Response `json:",inline"`

	Data *Data `json:"Data,omitempty"`
}

// Data 数据
// swagger:model JobListData
type Data struct {
	Jobs  []*schema.JobInfo `json:"Jobs"`
	Total int64             `json:"Total"`
}
