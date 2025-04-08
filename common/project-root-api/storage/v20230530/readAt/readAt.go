package readAt

import "github.com/yuansuan/ticp/common/project-root-api/pkg/xhttp"

type Request struct {
	Path       string                 `json:"Path" xquery:"Path" form:"Path" `
	Offset     int64                  `json:"Offset" xquery:"Offset" form:"Offset"`
	Length     int64                  `json:"Length" xquery:"Length" form:"Length"`
	Compressor string                 `json:"Compressor" xquery:"Compressor" form:"Compressor"`
	Resolver   xhttp.ResponseResolver `json:"-" xquery:"-" form:"-"`
}

type Response struct {
	Data []byte
}
