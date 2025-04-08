package writeAt

import (
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"io"
)

type Request struct {
	Path       string    `json:"Path" xquery:"Path" form:"Path" `
	Offset     int64     `json:"Offset" xquery:"Offset" form:"Offset"`
	Length     int64     `json:"Length" xquery:"Length" form:"Length"`
	Compressor string    `json:"Compressor" xquery:"Compressor" form:"Compressor"`
	Data       io.Reader `json:"—" xquery:"-" form:"-"`
}

type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"Data,omitempty"`
}

type Data struct {
	TransferSize int64 `json:"TransferSize"`
}
