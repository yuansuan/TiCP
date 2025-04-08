package init

import (
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// swagger:model storageUploadInitRequest
type Request struct {
	// 路径
	Path string `json:"Path"`
	// 文件大小
	Size int64 `json:"Size"`
	// 是否覆盖
	Overwrite bool `json:"Overwrite"`
}

// swagger:model storageUploadInitResponse
type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"Data,omitempty"`
}

// swagger:model storageUploadInitData
type Data struct {
	UploadID string `json:"UploadID"`
}
