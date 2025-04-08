package jobsnapshotget

import (
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
// swagger:model JobSnapshotGetRequest
type Request struct {
	JobID string `uri:"JobID" binding:"required"`                // 作业ID
	Path  string `query:"Path" xquery:"Path" binding:"required"` // 云图位置, 即list snapshot返回的数据中的值
}

// Response 返回
// swagger:model JobSnapshotGetResponse
type Response struct {
	schema.Response `json:",inline"`

	Data *Data `json:"Data,omitempty"`
}

// Data 数据
// swagger:model JobSnapshotGetData
type Data string
