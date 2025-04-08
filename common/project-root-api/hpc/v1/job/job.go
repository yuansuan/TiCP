package job

import (
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

type SystemPostRequest struct {
	IdempotentID string `uri:"IdempotentID"`
	// example "image:ubuntu:18.04"
	// example "local:appPath"
	Application     string                             `json:"Application,omitempty"`
	Environment     map[string]string                  `json:"Environment,omitempty"`
	Command         string                             `json:"Command,omitempty"`
	Override        v20230530.JobInHPCOverride         `json:"Override,omitempty"`
	Queue           string                             `json:"Queue,omitempty"`
	Resource        v20230530.JobInHPCResource         `json:"Resource,omitempty"`
	Inputs          []v20230530.JobInHPCInputStorage   `json:"Inputs,omitempty"`
	Output          *v20230530.JobInHPCOutputStorage   `json:"Output,omitempty"`
	CustomStateRule *v20230530.JobInHPCCustomStateRule `json:"CustomStateRule,omitempty"`
	// {"-o": "stdout.log", "-e", "stderr.log", "-l": "select=xxx"}
	JobSchedulerSubmitFlags map[string]string `json:"JobSchedulerSubmitFlags,omitempty"`
}

type SystemPostResponse struct {
	v20230530.Response

	Data *v20230530.JobInHPC `json:"Data,omitempty"`
}

type SystemGetRequest struct {
	JobID string `uri:"JobID"`
}

type SystemGetResponse struct {
	v20230530.Response

	Data *v20230530.JobInHPC `json:"Data,omitempty"`
}

type SystemCancelRequest struct {
	JobID string `uri:"JobID"`
}

type SystemCancelResponse struct {
	v20230530.Response

	Data *v20230530.JobInHPC `json:"Data,omitempty"`
}

type SystemDeleteRequest struct {
	JobID string `uri:"JobID"`
}

type SystemDeleteResponse struct {
	v20230530.Response
}

type SystemListRequest struct {
	PageOffset int      `query:"PageOffset"`
	PageSize   int      `query:"PageSize"`
	Status     string   `query:"Status"`
	JobIDs     []string `query:"JobIDs"`
}

type SystemListResponse struct {
	v20230530.Response

	Data *SystemListResponseData `json:"Data,omitempty"`
}

type SystemListResponseData struct {
	Jobs []*v20230530.JobInHPC `json:"Jobs,omitempty"`

	Offset int `json:"Offset"`
	Size   int `json:"Size"`
	Total  int `json:"Total"`
}
