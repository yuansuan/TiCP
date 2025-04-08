package batchDownload

import "github.com/yuansuan/ticp/common/project-root-api/pkg/xhttp"

type Request struct {
	Paths      []string               `form:"Paths" json:"Paths" xquery:"Paths"`
	BasePath   string                 `form:"BasePath" json:"BasePath" xquery:"BasePath"`
	FileName   string                 `form:"FileName" json:"FileName" xquery:"FileName"`
	IsCompress bool                   `form:"IsCompress" json:"IsCompress" xquery:"IsCompress"`
	Resolver   xhttp.ResponseResolver `form:"-" json:"-" xquery:"-"`
}

type Response struct {
	FileType string
	FileSize int64
}
