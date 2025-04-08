package file

import (
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// swagger:model storageUploadFileRequest
type Request struct {
	// 路径
	Path string `json:"Path" form:"Path" `
	// 文件数据
	Content []byte `json:"Content" form:"Content" `
	// 是否覆盖
	Overwrite bool `json:"Overwrite" form:"Overwrite"`
}

// swagger:model storageUploadFileResponse
type Response struct {
	v20230530.Response `json:",inline"`
}
