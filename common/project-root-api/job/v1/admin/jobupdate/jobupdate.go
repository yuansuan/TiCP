package jobupdate

import (
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
// swagger:model AdminJobUpdateRequest
type Request struct {
	JobID         string `uri:"JobID" binding:"required"` // 作业ID
	FileSyncState string `json:"FileSyncState"`           // 作业文件同步状态
}

// Response 返回
// swagger:model AdminJobUpdateResponse
type Response struct {
	schema.Response `json:",inline"`
}
