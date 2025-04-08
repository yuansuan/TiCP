package mkdir

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

// swagger:model storageMkdirRequest
type Request struct {
	// 路径
	Path string `json:"Path"`
	//是否忽略已存在
	IgnoreExist bool `json:"IgnoreExist"`
}

// swagger:model storageMkdirResponse
type Response struct {
	v20230530.Response `json:",inline"`
}
