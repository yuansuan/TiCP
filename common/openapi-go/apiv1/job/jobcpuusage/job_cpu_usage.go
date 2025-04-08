package jobcpuusage

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	jcu "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcpuusage"
)

type API func(options ...Option) (*jcu.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*jcu.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/api/jobs/" + req.JobID + "/cpuusage"))

		ret := new(jcu.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*jcu.Request, error) {
	req := &jcu.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *jcu.Request) error

func (api API) JobId(JobId string) Option {
	return func(req *jcu.Request) error {
		req.JobID = JobId
		return nil
	}
}
