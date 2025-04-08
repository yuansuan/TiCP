package status

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	CompressID string `form:"CompressID" json:"CompressID" query:"CompressID"`
}

type Response struct {
	v20230530.Response `json:",inline"`
	Data               *Data `json:"Data,omitempty"`
}

type Data struct {
	IsFinished bool   `json:"IsFinished"`
	Status     string `json:"Status"`
}
