package jobbatchget

import (
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
type Request struct {
	JobIDs []string `json:"JobIDs" binding:"required"` //作业ID
}

// Response 返回
type Response struct {
	schema.Response `json:",inline"`

	Data []*Data `json:"Data,omitempty"`
}

// Data 数据
type Data struct {
	schema.JobInfo `json:",inline"`
}
