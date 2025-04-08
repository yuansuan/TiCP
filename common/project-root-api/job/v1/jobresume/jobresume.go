package jobresume

import schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

// Request 请求
// swagger:model JobResumeRequest
type Request struct {
	JobID string `uri:"JobID" binding:"required"` //作业ID
}

type Response struct {
	schema.Response `json:",inline"`
}
