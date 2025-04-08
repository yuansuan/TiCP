package api

import schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

// ListSharedDirectoryRequest 获取共享目录列表
type ListSharedDirectoryRequest struct {
	PathPrefix string `json:"PathPrefix" query:"PathPrefix" form:"PathPrefix"`
}

// ListSharedDirectoryResponse 获取共享目录列表
type ListSharedDirectoryResponse struct {
	schema.Response `json:",inline"`
	Data            *Data `json:"Data,omitempty"`
}
