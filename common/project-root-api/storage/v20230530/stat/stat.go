package stat

import (
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// swagger:parameters storageStatRequest
type Request struct {
	// 路径
	// in: query
	Path string `form:"Path" query:"Path" json:"Path"`
}

// swagger:model storageStatResponse
type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"Data,omitempty"`
}

// swagger:model storageStatData
type Data struct {
	File *v20230530.FileInfo `json:"File"`
}
