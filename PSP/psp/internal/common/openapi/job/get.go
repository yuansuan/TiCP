package job

import (
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

// AdminGetJob 管理员获取作业信息
func AdminGetJob(api *openapi.OpenAPI, jobID string) (*schema.AdminJobInfo, error) {
	if strutil.IsEmpty(jobID) {
		return nil, ErrJobIDEmpty
	}

	option := api.Client.Job.AdminJobGet.JobId(jobID)
	resp, err := api.Client.Job.AdminJobGet(option)
	if err != nil {
		return nil, errors.Wrap(err, "openapi get admin job info err")
	}

	if resp != nil {
		logging.Default().Debugf("openapi admin get job request id: [%v], jobID: [%v]", resp.RequestID, jobID)
	}

	if resp.Data != nil {
		return &resp.Data.AdminJobInfo, err
	} else {
		return nil, ErrJobIsNotExist
	}
}

// GetJob 获取作业信息
func GetJob(api *openapi.OpenAPI, jobID string) (*schema.JobInfo, error) {
	if strutil.IsEmpty(jobID) {
		return nil, ErrJobIDEmpty
	}

	option := api.Client.Job.JobGet.JobId(jobID)
	resp, err := api.Client.Job.JobGet(option)
	if err != nil {
		return nil, errors.Wrap(err, "openapi get job info err")
	}

	if resp != nil {
		logging.Default().Debugf("openapi get job request id: [%v], jobID: [%v]", resp.RequestID, jobID)
	}

	if resp.Data != nil {
		return &resp.Data.JobInfo, err
	} else {
		return nil, ErrJobIsNotExist
	}
}
