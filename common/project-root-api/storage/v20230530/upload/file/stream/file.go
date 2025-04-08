package file

import (
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"io"
)

// swagger:model storageUploadFileRequest
type Request struct {
	// 路径
	Path string `json:"Path" form:"Path" `
	// 文件数据
	Content io.Reader `json:"Content" form:"Content" `
	// 是否覆盖
	Overwrite bool `json:"Overwrite" form:"Overwrite"`
	// 默认不压缩，目前只支持zip格式
	CompressType string `json:"CompressType" form:"CompressType"`
}

// swagger:model storageUploadFileResponse
type Response struct {
	v20230530.Response `json:",inline"`
}
