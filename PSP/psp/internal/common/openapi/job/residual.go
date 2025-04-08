package job

import (
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

func AdminGetResidual(api *openapi.OpenAPI, jobID string) (*schema.Residual, error) {
	if strutil.IsEmpty(jobID) {
		return nil, ErrJobIDEmpty
	}

	option := api.Client.Job.AdminJobGetResidual.JobId(jobID)
	resp, err := api.Client.Job.AdminJobGetResidual(option)
	if err != nil {
		return nil, errors.Wrap(err, "openapi get admin residual info err")
	}

	if resp != nil {
		logging.Default().Debugf("openapi admin get residual request id: [%v], jobID: [%v]", resp.RequestID, jobID)
	}

	if resp.Data != nil {
		return &resp.Data.Residual, err
	} else {
		return nil, ErrResidualIsNotExist
	}
}

func GetResidual(api *openapi.OpenAPI, jobID string) (*schema.Residual, error) {
	if strutil.IsEmpty(jobID) {
		return nil, ErrJobIDEmpty
	}

	option := api.Client.Job.JobGetResidual.JobId(jobID)
	resp, err := api.Client.Job.JobGetResidual(option)
	if err != nil {
		return nil, errors.Wrap(err, "openapi get residual info err")
	}

	if resp != nil {
		logging.Default().Debugf("openapi get residual request id: [%v], jobID: [%v]", resp.RequestID, jobID)
	}

	if resp.Data != nil {
		return &resp.Data.Residual, err
	} else {
		return nil, ErrResidualIsNotExist
	}
}
