package jobdelete

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	jd "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobdelete"
)

type API func(options ...Option) (*jd.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*jd.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			Method(http.MethodDelete).
			URI("/api/jobs/" + req.JobID))

		ret := new(jd.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*jd.Request, error) {
	req := &jd.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *jd.Request) error

func (api API) JobId(JobId string) Option {
	return func(req *jd.Request) error {
		req.JobID = JobId
		return nil
	}
}
