package jobresidual

import (
	residualget "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobresidual"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
// swagger:model JobResidualRequest
type Request struct {
	residualget.Request `json:",inline"`
}

// Response 返回
// swagger:model JobResidualResponse
type Response struct {
	schema.Response `json:",inline"`

	Data *Data `json:"Data,omitempty"`
}

// Data 数据
// swagger:model JobResidualData
type Data struct {
	schema.Residual `json:",inline"`
}
