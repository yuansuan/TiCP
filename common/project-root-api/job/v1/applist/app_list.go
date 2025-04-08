package applist

import (
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

type Request struct {
}

type Response struct {
	schema.Response `json:",inline"`

	Data []*schema.Application `json:"Data,omitempty"`
}
