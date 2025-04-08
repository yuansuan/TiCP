package resource

import (
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

type SystemGetRequest struct {
	Queue string `query:"Queue"`
}

type SystemGetResponse struct {
	v20230530.Response

	// Data map[]*v20230530.Resource `json:"Data,omitempty"`
	Data map[string]*v20230530.Resource `json:"Data,omitempty"`
}
