package start

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	Path string `form:"Path" json:"Path"`
}

type Response struct {
	v20230530.Response `json:",inline"`
	Data               *Data `json:"Data,omitempty"`
}

type Data struct {
	DirectoryUsageTaskID string `json:"DirectoryUsageTaskID"`
}
