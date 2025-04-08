package truncate

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

// swagger:model storageTruncateRequest
type Request struct {
	// 路径
	Path string `json:"Path"`
	// 文件大小
	Size int64 `json:"Size"`
	// 文件不存在时是否创建文件
	CreateIfNotExists bool `json:"CreateIfNotExists"`
}

// swagger:model storageTruncateResponse
type Response struct {
	v20230530.Response `json:",inline"`
}
