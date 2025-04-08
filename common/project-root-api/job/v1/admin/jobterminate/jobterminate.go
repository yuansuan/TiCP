package jobterminate

import (
	terminate "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobterminate"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
// swagger:model AdminJobTerminateRequest
type Request struct {
	terminate.Request `json:",inline"`
}

// Response 响应
type Response struct {
	schema.Response `json:",inline"`
}
