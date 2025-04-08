package quotadelete

import (
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
type Request struct {
	AppID  string `uri:"AppID"`
	UserID string `json:"UserID" binding:"required"`
}

// Response 响应
type Response struct {
	schema.Response `json:",inline"`
}
