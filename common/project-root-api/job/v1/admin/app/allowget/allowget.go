package allowget

import (
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
type Request struct {
	AppID string `uri:"AppID"`
}

// Response 响应
type Response struct {
	schema.Response `json:",inline"`

	Data *schema.ApplicationAllow `json:"Data,omitempty"`
}
