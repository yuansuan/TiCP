package jobterminate

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	jt "github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobterminate"
)

type API func(options ...Option) (*jt.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*jt.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			Method(http.MethodPatch).
			URI("/admin/jobs/" + req.JobID + "/terminate"))

		ret := new(jt.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*jt.Request, error) {
	req := &jt.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *jt.Request) error

func (api API) JobId(JobId string) Option {
	return func(req *jt.Request) error {
		req.JobID = JobId
		return nil
	}
}
