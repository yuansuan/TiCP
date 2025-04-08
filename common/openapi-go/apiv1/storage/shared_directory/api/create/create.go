package create

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/shared_directory/api"
)

// API 创建共享目录
type API func(options ...Option) (*api.CreateSharedDirectoryResponse, error)

// New 创建共享目录
func New(hc *xhttp.Client) API {
	return func(options ...Option) (*api.CreateSharedDirectoryResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/api/storage/sharedDirectorys").
			Json(req))

		ret := new(api.CreateSharedDirectoryResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

// NewRequest 创建共享目录
func NewRequest(options []Option) (*api.CreateSharedDirectoryRequest, error) {
	req := &api.CreateSharedDirectoryRequest{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

// Option 参数选项
type Option func(req *api.CreateSharedDirectoryRequest) error

// Paths 指定路径
func (a API) Paths(paths []string) Option {
	return func(req *api.CreateSharedDirectoryRequest) error {
		req.Paths = paths
		return nil
	}
}

// IgnoreExisting 是否忽略已存在的共享目录
func (a API) IgnoreExisting(ignoreExisting bool) Option {
	return func(req *api.CreateSharedDirectoryRequest) error {
		req.IgnoreExisting = ignoreExisting
		return nil
	}
}
