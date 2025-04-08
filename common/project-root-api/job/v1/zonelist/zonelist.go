package zonelist

import (
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
// swagger:model ZoneListRequest
type Request struct{}

// Response 返回
// swagger:model ZoneListResponse
type Response struct {
	schema.Response `json:",inline"`

	Data *Data `json:"Data,omitempty"`
}

// Data 数据
// swagger:model ZoneListData
type Data struct {
	schema.Zones `json:"Zones"`
}
