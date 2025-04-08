package list

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	app "github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/list"
)

// API 查询App列表API
type API func(options ...Option) (*app.Response, error)

// New 创建查询App列表API
func New(hc *xhttp.Client) API {
	return func(options ...Option) (*app.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		rb := xhttp.NewRequestBuilder().URI("/admin/apps")
		if req.AllowUserID != "" {
			rb.AddQuery("AllowUserID", req.AllowUserID)
		}

		resolver := hc.Prepare(rb)

		ret := new(app.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

// NewRequest 创建请求
func NewRequest(options []Option) (*app.Request, error) {
	req := &app.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

// Option 请求选项
type Option func(req *app.Request) error

// AllowUserID 查询允许的用户ID
func (api API) AllowUserID(allowUserID string) Option {
	return func(req *app.Request) error {
		req.AllowUserID = allowUserID
		return nil
	}
}
