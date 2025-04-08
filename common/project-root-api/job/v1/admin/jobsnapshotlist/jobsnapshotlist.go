package jobsnapshotlist

import (
	jobsnapshotlist "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobsnapshotlist"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
// swagger:model JobSnapshotListRequest
type Request struct {
	jobsnapshotlist.Request `json:",inline"`
}

// Response 返回
// swagger:model JobSnapshotListResponse
type Response struct {
	schema.Response `json:",inline"`

	Data *Data `json:"Data,omitempty"`
}

// Data 数据
// swagger:model JobSnapshotListData
type Data map[string][]string // key:云图集名 value:对应云图集文件名列表
