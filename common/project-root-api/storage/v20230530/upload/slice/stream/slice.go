package slice

import (
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"io"
)

// swagger:model storageUploadSliceRequest
type Request struct {
	// 文件上传ID
	UploadID string `json:"UploadID"`
	// 写入offset
	Offset int64 `json:"Offset"`
	// 写入长度
	Length int64 `json:"Length"`
	// 写入数据
	Slice io.Reader `json:"Slice"`
	// 压缩方式
	Compressor string `json:"Compressor"`
}

// swagger:model storageUploadSliceResponse
type Response struct {
	v20230530.Response `json:",inline"`
}
