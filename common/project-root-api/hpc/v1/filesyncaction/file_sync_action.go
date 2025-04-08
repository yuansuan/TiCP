package filesyncaction

import (
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

type SystemPauseFileSyncRequest struct {
	JobID string `uri:"JobID"`
}

type SystemPauseFileSyncResponse struct {
	v20230530.Response

	Data *v20230530.JobInHPC `json:"Data,omitempty"`
}

type SystemResumeFileSyncRequest struct {
	JobID string `uri:"JobID"`
}

type SystemResumeFileSyncResponse struct {
	v20230530.Response

	Data *v20230530.JobInHPC `json:"Data,omitempty"`
}
