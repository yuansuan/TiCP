package jobbatchget

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	jg "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobbatchget"
)

// API 批量获取作业api
type API func(options ...Option) (*jg.Response, error)

// New 新建
func New(hc *xhttp.Client) API {
	return func(options ...Option) (*jg.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/api/jobs/batch").
			Json(req))

		ret := new(jg.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

// NewRequest 新建请求
func NewRequest(options []Option) (*jg.Request, error) {
	req := &jg.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

// Option 选项
type Option func(req *jg.Request) error

// JobIds 作业IDs
func (api API) JobIds(JobIds []string) Option {
	return func(req *jg.Request) error {
		req.JobIDs = JobIds
		return nil
	}
}
