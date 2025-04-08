package status

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	DirectoryUsageTaskID string `form:"DirectoryUsageTaskID" json:"DirectoryUsageTaskID" query:"DirectoryUsageTaskID"`
}

type Response struct {
	v20230530.Response `json:",inline"`
	Data               *Data `json:"Data,omitempty"`
}

type Data struct {
	Status     string `json:"Status"`
	Size       int64  `json:"Size"`
	LogicSize  int64  `json:"LogicSize"`
	ErrMessage string `json:"ErrMessage"`
}
