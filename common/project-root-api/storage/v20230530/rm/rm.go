package rm

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

// swagger:model storageRmRequest
type Request struct {
	// 路径
	Path string `json:"Path"`
	//是否忽略不存在
	IgnoreNotExist bool `json:"IgnoreNotExist"`
}

// swagger:model storageRmResponse
type Response struct {
	v20230530.Response `json:",inline"`
}
