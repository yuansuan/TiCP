package job

import (
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	jobsnapshotgetadmin "github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobsnapshotget"
	jobsnapshotlistadmin "github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobsnapshotlist"
	jobsnapshotget "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobsnapshotget"
	jobsnapshotlist "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobsnapshotlist"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

func AdminListSnapshot(api *openapi.OpenAPI, jobID string) (*jobsnapshotlistadmin.Response, error) {
	if strutil.IsEmpty(jobID) {
		return nil, ErrJobIDEmpty
	}

	option := api.Client.Job.AdminJobListSnapshot.JobId(jobID)
	resp, err := api.Client.Job.AdminJobListSnapshot(option)
	if err != nil {
		return nil, errors.Wrap(err, "openapi admin list snapshot err")
	}

	if resp != nil {
		logging.Default().Debugf("openapi admin list snapshot request id: [%v], jobID: [%v]", resp.RequestID, jobID)
	}

	if resp.Data != nil {
		return resp, err
	} else {
		return nil, ErrSnapshotsIsNotExist
	}
}

func AdminGetSnapshot(api *openapi.OpenAPI, jobID, path string) (*jobsnapshotgetadmin.Response, error) {
	if strutil.IsEmpty(jobID) {
		return nil, ErrJobIDEmpty
	}
	if strutil.IsEmpty(path) {
		return nil, ErrSnapshotPathIsNotExist
	}

	resp, err := api.Client.Job.AdminJobGetSnapshot(
		api.Client.Job.AdminJobGetSnapshot.JobId(jobID),
		api.Client.Job.AdminJobGetSnapshot.Path(path),
	)
	if err != nil {
		return nil, errors.Wrap(err, "openapi admin get snapshot err")
	}

	if resp != nil {
		logging.Default().Debugf("openapi admin get snapshot request id: [%v], jobID: [%v], path: [%v]", resp.RequestID, jobID, path)
	}

	if resp.Data != nil {
		return resp, err
	} else {
		return nil, ErrSnapshotIsNotExist
	}
}

func ListSnapshot(api *openapi.OpenAPI, jobID string) (*jobsnapshotlist.Response, error) {
	if strutil.IsEmpty(jobID) {
		return nil, ErrJobIDEmpty
	}

	option := api.Client.Job.JobListSnapshot.JobId(jobID)
	resp, err := api.Client.Job.JobListSnapshot(option)
	if err != nil {
		return nil, errors.Wrap(err, "openapi list snapshot err")
	}

	if resp != nil {
		logging.Default().Debugf("openapi list snapshot request id: [%v], jobID: [%v]", resp.RequestID, jobID)
	}

	if resp.Data != nil {
		return resp, err
	} else {
		return nil, ErrSnapshotsIsNotExist
	}
}

func GetSnapshot(api *openapi.OpenAPI, jobID, path string) (*jobsnapshotget.Response, error) {
	if strutil.IsEmpty(jobID) {
		return nil, ErrJobIDEmpty
	}
	if strutil.IsEmpty(path) {
		return nil, ErrSnapshotPathIsNotExist
	}

	resp, err := api.Client.Job.JobGetSnapshot(
		api.Client.Job.JobGetSnapshot.JobId(jobID),
		api.Client.Job.JobGetSnapshot.Path(path),
	)
	if err != nil {
		return nil, errors.Wrap(err, "openapi get snapshot err")
	}

	if resp != nil {
		logging.Default().Debugf("openapi get snapshot request id: [%v], jobID: [%v], path: [%v]", resp.RequestID, jobID, path)
	}

	if resp.Data != nil {
		return resp, err
	} else {
		return nil, ErrSnapshotIsNotExist
	}
}
