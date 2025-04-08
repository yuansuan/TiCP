package admin

import (
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

type ListOperationLogRequest struct {
	// 用户 ID
	UserID string `form:"UserID" json:"UserID" query:"UserID"`
	// 文件名
	FileName string `form:"FileName" json:"FileName" query:"FileName"`
	// 文件类型, 可选值: file-普通文件, folder-文件夹, batch-批量操作
	FileTypes string `form:"FileTypes" json:"FileTypes" query:"FileTypes"`
	// 操作类型, 可选值: upload-上传, download-下载, delete-删除, move-移动, mkdir-添加文件夹, copy-拷贝, copy_range-指定范围拷贝,compress-压缩, create-创建, link-链接, read_at-读, write_at-写
	OperationTypes string `form:"OperationTypes" json:"OperationTypes" query:"OperationTypes"`
	// 开始时间戳
	BeginTime int64 `form:"BeginTime" json:"BeginTime" query:"BeginTime"`
	// 结束时间戳
	EndTime int64 `form:"EndTime" json:"EndTime" query:"EndTime"`
	//分页偏移量
	PageOffset int64 `form:"PageOffset" json:"PageOffset" query:"PageOffset"`
	//分页大小
	PageSize int64 `form:"PageSize" json:"PageSize" query:"PageSize"`
}

type ListOperationLogResponse struct {
	v20230530.Response `json:",inline"`

	Data *ListOperationLogResponseData `json:"Data,omitempty"`
}

type ListOperationLogResponseData struct {
	OperationLog []*v20230530.OperationLog `json:"OperationLog"`
	NextMarker   int64                     `json:"NextMarker"`
	Total        int64                     `json:"Total"`
}
