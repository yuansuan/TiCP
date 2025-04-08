package jobpreschedule

import (
	"net/http"

	js "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobpreschedule"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

// API 提交作业api
type API func(options ...Option) (*js.Response, error)

// New 新建
func New(hc *xhttp.Client) API {
	return func(options ...Option) (*js.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/api/jobs/preschedule").
			Json(req))

		ret := new(js.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

// NewRequest 新建请求
func NewRequest(options []Option) (*js.Request, error) {
	req := &js.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

// Option 选项
type Option func(req *js.Request) error

// Zones 预选分区列表
func (api API) Zones(zones []string) Option {
	return func(req *js.Request) error {
		req.Zones = zones
		return nil
	}
}

// Params 预调度作业参数
func (api API) Params(params js.Params) Option {
	return func(req *js.Request) error {
		req.Params = params
		return nil
	}
}

func (api API) Shared(shared bool) Option {
	return func(req *js.Request) error {
		req.Shared = shared
		return nil
	}
}

func (api API) Fixed(fixed bool) Option {
	return func(req *js.Request) error {
		req.Fixed = fixed
		return nil
	}
}
