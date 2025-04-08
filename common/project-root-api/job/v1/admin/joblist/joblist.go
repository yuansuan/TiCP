package joblist

import (
	list "github.com/yuansuan/ticp/common/project-root-api/job/v1/joblist"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
// swagger:model AdminJobListRequest
type Request struct {
	list.Request   `json:",inline"`
	UserID         string `form:"UserID"`
	AppID          string `form:"AppID"`
	WithDelete     bool   `form:"WithDelete"`
	IsSystemFailed bool   `form:"IsSystemFailed"`
}

// Response 返回
// swagger:model AdminJobListResponse
type Response struct {
	schema.Response `json:",inline"`

	Data *Data `json:"Data,omitempty"`
}

// Data 数据
// swagger:model AdminJobListData
type Data struct {
	Jobs  []*schema.AdminJobInfo `json:"Jobs"`
	Total int64                  `json:"Total"`
}
