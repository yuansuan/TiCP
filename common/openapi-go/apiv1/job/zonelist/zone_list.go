package zonelist

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	zl "github.com/yuansuan/ticp/common/project-root-api/job/v1/zonelist"
)

// API 分区列表api
type API func(options ...Option) (*zl.Response, error)

// New 新建
func New(hc *xhttp.Client) API {
	return func(options ...Option) (*zl.Response, error) {

		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/api/zones"))

		ret := new(zl.Response)
		err := resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

// NewRequest 新建请求
func NewRequest(options []Option) (*zl.Request, error) {
	req := &zl.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

// Option 选项
type Option func(req *zl.Request) error
