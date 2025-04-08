package list

import (
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
type Request struct {
	AllowUserID string `query:"AllowUserID" xquery:"AllowUserID"`
}

// Response 响应
type Response struct {
	schema.Response `json:",inline"`

	Data *[]schema.Application `json:"Data,omitempty"`
}
