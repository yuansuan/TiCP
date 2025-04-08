package mv

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

// swagger:model storageMvRequest
type Request struct {
	// 源路径
	SrcPath string `json:"Src"`
	// 目标路径
	DestPath string `json:"Dest"`
}

// swagger:model storageMvResponse
type Response struct {
	v20230530.Response `json:",inline"`
}
