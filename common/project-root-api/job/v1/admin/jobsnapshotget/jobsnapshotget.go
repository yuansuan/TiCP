package jobsnapshotget

import (
	jobsnapshotget "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobsnapshotget"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
// swagger:model JobSnapshotGetRequest
type Request struct {
	jobsnapshotget.Request `json:",inline"`
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
