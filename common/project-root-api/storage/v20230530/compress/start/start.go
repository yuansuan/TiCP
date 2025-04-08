package start

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	Paths      []string `form:"Paths" json:"Paths"`
	TargetPath string   `form:"TargetPath" json:"TargetPath"`
	BasePath   string   `form:"BasePath" json:"BasePath" xquery:"BasePath"`
}

type Response struct {
	v20230530.Response `json:",inline"`
	Data               *Data `json:"Data,omitempty"`
}

type Data struct {
	CompressID string `json:"CompressID"`
	FileName   string `json:"FileName"`
}
