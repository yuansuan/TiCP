package copyRange

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	SrcPath    string `json:"SrcPath"`
	DestPath   string `json:"DestPath"`
	SrcOffset  int64  `json:"SrcOffset"`
	DestOffset int64  `json:"DestOffset"`
	Length     int64  `json:"Length"`
}

type Response struct {
	v20230530.Response `json:",inline"`
}
