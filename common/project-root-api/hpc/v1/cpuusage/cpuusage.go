package cpuusage

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type SystemGetRequest struct {
	JobID string `uri:"JobID"`
}

type SystemGetResponse struct {
	v20230530.Response

	Data *v20230530.CpuUsage `json:"Data,omitempty"`
}
