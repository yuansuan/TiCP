package jobretransmit

import schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

// Request 请求
// swagger:model jobRetransmitRequest
type Request struct {
	JobID string `uri:"JobID" binding:"required"` //作业ID
}

// Response 响应
// swagger:model JobDeleteResponse
type Response struct {
	schema.Response `json:",inline"`
}
