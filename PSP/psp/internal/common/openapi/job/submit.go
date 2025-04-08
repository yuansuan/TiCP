package job

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	adminjobcreate "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/jobcreate"
	job "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcreate"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

// SubmitParams 作业提交参数信息
type SubmitParams struct {
	Name            string
	Zone            string
	Queue           string
	Comment         string
	Params          job.Params
	SchedulerParams map[string]string
}

// SubmitJob 提交作业
func SubmitJob(api *openapi.OpenAPI, submitParams *SubmitParams) (string, error) {
	submitJob := api.Client.Job.AdminJobCreate
	options := []adminjobcreate.Option{
		submitJob.Name(submitParams.Name),
		submitJob.Zone(submitParams.Zone),
		submitJob.Queue(submitParams.Queue),
		submitJob.Comment(submitParams.Comment),
		submitJob.NoRound(true),
		submitJob.Params(submitParams.Params),
		submitJob.JobSchedulerSubmitFlags(submitParams.SchedulerParams),
	}

	resp, err := submitJob(options...)
	if err != nil {
		return "", err
	}

	if resp != nil {
		logging.Default().Debugf("openapi submit local job request id: [%v], req: [%+v]", resp.RequestID, submitParams)
	}

	if resp.Data != nil && !strutil.IsEmpty(resp.Data.JobID) {
		return resp.Data.JobID, err
	} else {
		return "", ErrJobIDEmpty
	}
}
