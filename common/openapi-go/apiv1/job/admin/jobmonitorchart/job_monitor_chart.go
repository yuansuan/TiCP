package jobmonitorchart

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	jm "github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobmonitorchart"
)

type API func(options ...Option) (*jm.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*jm.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/admin/jobs/" + req.JobID + "/monitorchart"))

		ret := new(jm.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*jm.Request, error) {
	req := &jm.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *jm.Request) error

func (api API) JobId(JobId string) Option {
	return func(req *jm.Request) error {
		req.JobID = JobId
		return nil
	}
}
