package complete

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/complete"
	"net/http"
)

type API func(options ...Option) (*complete.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*complete.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/system/storage/upload/complete").
			Json(req))

		ret := new(complete.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*complete.Request, error) {
	req := &complete.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *complete.Request) error

func (api API) Path(path string) Option {
	return func(req *complete.Request) error {
		req.Path = path
		return nil
	}
}

func (api API) UploadID(uploadID string) Option {
	return func(req *complete.Request) error {
		req.UploadID = uploadID
		return nil
	}
}
