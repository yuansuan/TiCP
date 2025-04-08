package list

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/quota/admin"
	"net/http"
)

type API func(options ...Option) (*admin.ListStorageQuotaResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*admin.ListStorageQuotaResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/admin/storage/storageQuota").
			AddQuery("PageOffset", utils.Stringify(req.PageOffset)).
			AddQuery("PageSize", utils.Stringify(req.PageSize)))

		ret := new(admin.ListStorageQuotaResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*admin.ListStorageQuotaRequest, error) {
	req := &admin.ListStorageQuotaRequest{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *admin.ListStorageQuotaRequest) error

func (api API) PageSize(pageSize int) Option {
	return func(req *admin.ListStorageQuotaRequest) error {
		req.PageSize = pageSize
		return nil
	}
}

func (api API) PageOffset(pageOffset int) Option {
	return func(req *admin.ListStorageQuotaRequest) error {
		req.PageOffset = pageOffset
		return nil
	}
}
