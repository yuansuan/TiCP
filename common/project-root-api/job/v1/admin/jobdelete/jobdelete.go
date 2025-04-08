package jobdelete

import (
	del "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobdelete"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
// swagger:model JobDeleteRequest
type Request struct {
	del.Request `json:",inline"`
}

// Response 响应
// swagger:model JobDeleteResponse
type Response struct {
	schema.Response `json:",inline"`
}
