package download

import "github.com/yuansuan/ticp/common/project-root-api/pkg/xhttp"

// swagger:parameters storageDownloadRequest
type Request struct {
	// 文件路径
	// in: query
	Path string `form:"Path" xquery:"Path" json:"Path"`
	// 文件分片range
	// in: header
	// see http range header Range: bytes={start}-{end}
	Range string `xquery:"Range" json:"Range" xheader:"Range"`

	Resolver xhttp.ResponseResolver
}

type Response struct {
	Filename string
	FileType string
	FileSize int64
	Data     []byte
}
