package link

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	SrcPath  string `json:"SrcPath"`
	DestPath string `json:"DestPath"`
}

// swagger:model storageMkdirResponse
type Response struct {
	v20230530.Response `json:",inline"`
}
