package job

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	apiJobList "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/joblist"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/joblist"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

// AdminListJobs 管理员分页获取作业列表
func AdminListJobs(api *openapi.OpenAPI, zone, jobState string, pageIndex, pageSize int64) (*joblist.Response, error) {
	if pageSize > xtype.MaxPageSize {
		return nil, ErrPageSizeInvalid
	}

	var options []apiJobList.Option
	listJobs := api.Client.Job.AdminJobList
	options = append(options, listJobs.Zone(zone))
	options = append(options, listJobs.JobState(jobState))
	options = append(options, listJobs.PageSize(pageSize))
	options = append(options, listJobs.PageOffset(pageIndex))

	response, err := listJobs(options...)

	if response != nil {
		logging.Default().Debugf("openapi admin list jobs request id: [%v], zone: [%v], jobState: [%v], pageIndex: [%v], "+
			"pageSize: [%v]", response.RequestID, zone, jobState, pageIndex, pageSize)
	}

	return response, err
}
