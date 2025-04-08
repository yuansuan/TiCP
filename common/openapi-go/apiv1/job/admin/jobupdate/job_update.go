package jobupdate

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	ju "github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobupdate"
)

type API func(options ...Option) (*ju.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*ju.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			Method(http.MethodPatch).
			URI("/admin/jobs/" + req.JobID + "/update").
			Json(req))
		ret := new(ju.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*ju.Request, error) {
	req := &ju.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *ju.Request) error

func (API) JobId(jobID string) Option {
	return func(req *ju.Request) error {
		req.JobID = jobID
		return nil
	}
}

func (API) FileSyncState(fileSyncState string) Option {
	return func(req *ju.Request) error {
		req.FileSyncState = fileSyncState
		return nil
	}
}
