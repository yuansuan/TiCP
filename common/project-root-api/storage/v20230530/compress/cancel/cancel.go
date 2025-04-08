package cancel

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	CompressID string `json:"CompressID"`
}

type Response struct {
	v20230530.Response `json:",inline"`
}
