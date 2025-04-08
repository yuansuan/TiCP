package job

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	admin "github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobterminate"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobterminate"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

// AdminTerminate 管理员终止作业
func AdminTerminate(api *openapi.OpenAPI, jobID string) (*admin.Response, error) {
	if strutil.IsEmpty(jobID) {
		return nil, ErrJobIDEmpty
	}

	terminateJob := api.Client.Job.AdminJobTerminate
	options := terminateJob.JobId(jobID)
	response, err := terminateJob(options)

	if response != nil {
		logging.Default().Debugf("openapi admin terminate request id: [%v], jobID: [%v]", response.RequestID, jobID)
	}

	return response, err
}

// Terminate 终止作业
func Terminate(api *openapi.OpenAPI, jobID string) (*jobterminate.Response, error) {
	if strutil.IsEmpty(jobID) {
		return nil, ErrJobIDEmpty
	}

	terminateJob := api.Client.Job.JobTerminate
	options := terminateJob.JobId(jobID)
	response, err := terminateJob(options)

	if response != nil {
		logging.Default().Debugf("openapi terminate request id: [%v], jobID: [%v]", response.RequestID, jobID)
	}

	return response, err
}
