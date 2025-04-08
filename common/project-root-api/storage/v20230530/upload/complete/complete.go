package complete

import (
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// swagger:model storageUploadCompleteRequest
type Request struct {
	// 路径
	Path string `json:"Path"`
	// 文件上传ID
	UploadID string `json:"UploadID"`
}

// swagger:model storageUploadCompleteResponse
type Response struct {
	v20230530.Response `json:",inline"`
}
