package jobsyncfilestate

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"net/http"

	jsfs "github.com/yuansuan/ticp/common/project-root-api/job/v1/system/jobsyncfilestate"
)

type API func(options ...Option) (*jsfs.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*jsfs.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			Method(http.MethodPatch).
			URI("/system/jobs/" + req.JobID + "/syncfile").Json(req))

		ret := new(jsfs.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*jsfs.Request, error) {
	req := &jsfs.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *jsfs.Request) error

func (api API) JobId(JobId string) Option {
	return func(req *jsfs.Request) error {
		req.JobID = JobId
		return nil
	}
}

func (api API) DownloadFinished(downloadFinished bool) Option {
	return func(req *jsfs.Request) error {
		req.DownloadFinished = downloadFinished
		return nil
	}
}

func (api API) DownloadFileSizeCurrent(downloadFileSizeCurrent int64) Option {
	return func(req *jsfs.Request) error {
		req.DownloadFileSizeCurrent = downloadFileSizeCurrent
		return nil
	}
}

func (api API) DownloadFileSizeTotal(downloadFileSizeTotal int64) Option {
	return func(req *jsfs.Request) error {
		req.DownloadFileSizeTotal = downloadFileSizeTotal
		return nil
	}
}

func (api API) DownloadFinishedTime(downloadFinishedTime string) Option {
	return func(req *jsfs.Request) error {
		req.DownloadFinishedTime = downloadFinishedTime
		return nil
	}
}

func (api API) TransmittingTime(transmittingTime string) Option {
	return func(req *jsfs.Request) error {
		req.TransmittingTime = transmittingTime
		return nil
	}
}

func (api API) FileSyncState(fileSyncState string) Option {
	return func(req *jsfs.Request) error {
		req.FileSyncState = fileSyncState
		return nil
	}
}
