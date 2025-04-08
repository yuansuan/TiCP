package list

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/shared_directory/api"
)

// API 查询共享目录
type API func(options ...Option) (*api.ListSharedDirectoryResponse, error)

// New 查询共享目录
func New(hc *xhttp.Client) API {
	return func(options ...Option) (*api.ListSharedDirectoryResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/api/storage/sharedDirectorys").
			AddQuery("PathPrefix", req.PathPrefix))

		ret := new(api.ListSharedDirectoryResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

// NewRequest 查询共享目录
func NewRequest(options []Option) (*api.ListSharedDirectoryRequest, error) {
	req := &api.ListSharedDirectoryRequest{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

// Option 参数选项
type Option func(req *api.ListSharedDirectoryRequest) error

// PathPrefix 指定路径前缀
func (a API) PathPrefix(pathPrefix string) Option {
	return func(req *api.ListSharedDirectoryRequest) error {
		req.PathPrefix = pathPrefix
		return nil
	}
}
