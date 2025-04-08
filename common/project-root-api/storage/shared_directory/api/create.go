package api

import schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

// CreateSharedDirectoryRequest 创建共享目录
type CreateSharedDirectoryRequest struct {
	IgnoreExisting bool     `json:"IgnoreExisting" form:"IgnoreExisting" query:"IgnoreExisting"`
	Paths          []string `json:"Paths" form:"Paths" query:"Paths"`
}

// CreateSharedDirectoryResponse 创建共享目录
type CreateSharedDirectoryResponse struct {
	schema.Response `json:",inline"`

	Data *Data `json:"Data,omitempty"`
}

// Data 数据
type Data []*schema.SharedDirectory
