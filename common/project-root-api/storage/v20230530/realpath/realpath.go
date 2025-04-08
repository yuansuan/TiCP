package realpath

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	RelativePath string `json:"RelativePath"`
}

// swagger:model storageMkdirResponse
type Response struct {
	v20230530.Response `json:",inline"`
	Data               *Data `json:"Data,omitempty"`
}

type Data struct {
	RealPath string `json:"RealPath"`
}
