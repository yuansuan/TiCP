package cancel

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	DirectoryUsageTaskID string `json:"DirectoryUsageTaskID"`
}

type Response struct {
	v20230530.Response `json:",inline"`
}
