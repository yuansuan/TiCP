package delete

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/shared_directory/api"
)

// API 删除共享目录
type API func(options ...Option) (*api.DeleteSharedDirectoryResponse, error)

// New 删除共享目录
func New(hc *xhttp.Client) API {
	return func(options ...Option) (*api.DeleteSharedDirectoryResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewDeleteRequestBuilder().
			URI("/api/storage/sharedDirectorys").
			Json(req))

		ret := new(api.DeleteSharedDirectoryResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

// NewRequest 删除共享目录
func NewRequest(options []Option) (*api.DeleteSharedDirectoryRequest, error) {
	req := &api.DeleteSharedDirectoryRequest{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

// Option 参数选项
type Option func(req *api.DeleteSharedDirectoryRequest) error

// Paths 指定路径
// 例如：["/YSID/dir1", "/YSID/dir2"]
func (a API) Paths(paths []string) Option {
	return func(req *api.DeleteSharedDirectoryRequest) error {
		req.Paths = paths
		return nil
	}
}

// IgnoreNonexistent 是否忽略不存在的共享目录
func (a API) IgnoreNonexistent(ignoreNonexistent bool) Option {
	return func(req *api.DeleteSharedDirectoryRequest) error {
		req.IgnoreNonexistent = ignoreNonexistent
		return nil
	}
}
