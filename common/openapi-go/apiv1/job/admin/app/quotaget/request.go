package quotaget

import (
	"errors"
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	appquota "github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/quotaget"
)

// API 获取应用配额
type API func(options ...Option) (*appquota.Response, error)

// New 创建API
func New(hc *xhttp.Client) API {
	return func(options ...Option) (*appquota.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		if req.AppID == "" {
			return nil, errors.New("AppID is required")
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/admin/apps/"+req.AppID+"/quota").
			AddQuery("UserID", req.UserID))

		ret := new(appquota.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

// NewRequest 创建请求
func NewRequest(options []Option) (*appquota.Request, error) {
	req := &appquota.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

// Option 选项
type Option func(req *appquota.Request) error

// AppID 应用ID
func (api API) AppID(AppID string) Option {
	return func(req *appquota.Request) error {
		req.AppID = AppID
		return nil
	}
}

// UserID 用户ID
func (api API) UserID(UserID string) Option {
	return func(req *appquota.Request) error {
		req.UserID = UserID
		return nil
	}
}
