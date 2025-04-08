package api

import schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

// DeleteSharedDirectoryRequest 删除共享目录
type DeleteSharedDirectoryRequest struct {
	IgnoreNonexistent bool     `json:"IgnoreNonexistent" form:"IgnoreNonexistent" query:"IgnoreNonexistent"`
	Paths             []string `json:"Paths" form:"Paths" query:"Paths"`
}

// DeleteSharedDirectoryResponse 删除共享目录
type DeleteSharedDirectoryResponse struct {
	schema.Response `json:",inline"`
}
