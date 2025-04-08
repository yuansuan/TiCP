package update

import (
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/add"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

type Status string

const (
	Published   Status = "published"
	Unpublished Status = "unpublished"
)

// Request 请求
type Request struct {
	AppID string `json:"AppID,omitempty" xquery:"AppID" form:"AppID" `
	add.Request
	PublishStatus Status `json:"PublishStatus,omitempty" xquery:"PublishStatus" form:"PublishStatus" `
}

// Response 响应
type Response struct {
	schema.Response `json:",inline"`
}
