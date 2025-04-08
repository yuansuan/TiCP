package allowget

import (
	"errors"
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	appallow "github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/allowget"
)

// API 获取应用配额
type API func(options ...Option) (*appallow.Response, error)

// New 创建API
func New(hc *xhttp.Client) API {
	return func(options ...Option) (*appallow.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		if req.AppID == "" {
			return nil, errors.New("AppID is required")
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/admin/apps/" + req.AppID + "/allow"))

		ret := new(appallow.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

// NewRequest 创建请求
func NewRequest(options []Option) (*appallow.Request, error) {
	req := &appallow.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

// Option 选项
type Option func(req *appallow.Request) error

// AppID 应用ID
func (api API) AppID(AppID string) Option {
	return func(req *appallow.Request) error {
		req.AppID = AppID
		return nil
	}
}
