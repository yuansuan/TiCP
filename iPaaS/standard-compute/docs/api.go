package docs

import (
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/job"
)

// --------------------------------------------------------------------------

// swagger:route POST /system/jobs Standard-Compute postJobsRequest
// 提交作业
// responses:
//   200: postJobsSystemResponse
//   400: errorResponse
//   401: errorResponse
//   500: errorResponse

// swagger:parameters postJobsRequest
type postJobsRequest struct {
	// This text will appear as description of your request body.
	// in:body
	Body job.SystemPostRequest
}

// swagger:response postJobsResponse
type postJobsResponse struct {
	// in:body
	Body *job.SystemPostResponse
}

// --------------------------------------------------------------------------

// swagger:route GET /system/jobs/:job_id Standard-Compute getJobRequest
// 查询作业
// responses:
//   200: getJobResponse
//   400: errorResponse
//   401: errorResponse
//   500: errorResponse

// swagger:parameters getJobRequest
type getJobRequest struct {
	job.SystemGetRequest
}

// swagger:response getJobResponse
type getJobResponse struct {
	// in:body
	Body *job.SystemGetResponse
}
