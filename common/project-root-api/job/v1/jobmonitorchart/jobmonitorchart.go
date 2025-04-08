package jobmonitorchart

import (
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
// swagger:model JobMonitorChartRequest
type Request struct {
	JobID string `uri:"JobID" binding:"required"` //作业ID
}

// Response 返回
// swagger:model JobMonitorChartResponse
type Response struct {
	schema.Response `json:",inline"`

	Data *Data `json:"Data,omitempty"`
}

// Data 数据
// swagger:model JobMonitorChartData
type Data []*schema.MonitorChart
