package file

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/file"
	"net/http"
)

type API func(options ...Option) (*file.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*file.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/api/storage/upload/file").
			AddQuery("Path", utils.Stringify(req.Path)).
			AddQuery("Overwrite", utils.Stringify(req.Overwrite)).
			BytesBody(req.Content))

		ret := new(file.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*file.Request, error) {
	req := &file.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *file.Request) error

func (api API) Path(path string) Option {
	return func(req *file.Request) error {
		req.Path = path
		return nil
	}
}

func (api API) Content(Content []byte) Option {
	return func(req *file.Request) error {
		req.Content = Content
		return nil
	}
}

func (api API) Overwrite(overwrite bool) Option {
	return func(req *file.Request) error {
		req.Overwrite = overwrite
		return nil
	}
}
