package delete

import (
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

type Request struct {
	AppID string `uri:"AppID" binding:"required"`
}

type Response struct {
	schema.Response `json:",inline"`
}
