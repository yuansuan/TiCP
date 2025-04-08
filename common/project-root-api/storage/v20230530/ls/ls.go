package ls

import (
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// swagger:parameters storageLsRequest
type Request struct {
	// 文件路径
	// in:query
	Path string `form:"Path" json:"Path" query:"Path"`
	// 分页偏移量
	// in:query
	PageOffset int64 `form:"PageOffset" json:"PageOffset" query:"PageOffset"`
	// 分页大小
	// in:query
	PageSize int64 `form:"PageSize" json:"PageSize" query:"PageSize"`
	// 过滤正则表达式 匹配的文件不会返回
	// in:query
	FilterRegexp string `form:"FilterRegexp" json:"FilterRegexp" query:"FilterRegexp"`
	// 过滤正则表达式list 匹配的文件不会返回
	// in:query
	FilterRegexpList []string `form:"FilterRegexpList" json:"FilterRegexpList" query:"FilterRegexpList"`
}

// swagger:model storageLsResponse
type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"Data,omitempty"`
}

// swagger:model storageLsData
type Data struct {
	Files      []*v20230530.FileInfo `json:"Files"`
	NextMarker int64                 `json:"NextMarker"`
	Total      int64                 `json:"Total"`
}
