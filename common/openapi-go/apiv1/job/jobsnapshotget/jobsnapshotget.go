package jobsnapshotget

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	jg "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobsnapshotget"
)

type API func(options ...Option) (*jg.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*jg.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}

		rb := xhttp.NewRequestBuilder().
			URI("/api/jobs/"+req.JobID+"/snapshots/img").
			AddQuery("Path", req.Path)

		resolver := hc.Prepare(rb)

		ret := new(jg.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*jg.Request, error) {
	req := &jg.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *jg.Request) error

func (api API) JobId(JobId string) Option {
	return func(req *jg.Request) error {
		req.JobID = JobId
		return nil
	}
}

func (api API) Path(path string) Option {
	return func(req *jg.Request) error {
		req.Path = path
		return nil
	}
}
