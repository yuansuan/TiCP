package readAt

import (
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	xhttp2 "github.com/yuansuan/ticp/common/project-root-api/pkg/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/readAt"
	"io"
	"net/http"
)

type API func(options ...Option) (*readAt.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*readAt.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}

		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/api/storage/readAt").
			Ref(req))

		ret := new(readAt.Response)
		f := func(resp *http.Response) error {
			return defaultResolver(resp, ret)
		}
		if req.Resolver == nil {
			req.Resolver = f
		}

		err = resolver.Resolve(xhttp.ResponseResolver(req.Resolver))

		return ret, err
	}
}

func NewRequest(options []Option) (*readAt.Request, error) {
	req := &readAt.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *readAt.Request) error

func (api API) Path(path string) Option {
	return func(req *readAt.Request) error {
		req.Path = path
		return nil
	}
}

func (api API) Offset(Offset int64) Option {
	return func(req *readAt.Request) error {
		req.Offset = Offset
		return nil
	}
}

func (api API) Length(length int64) Option {
	return func(req *readAt.Request) error {
		req.Length = length
		return nil
	}
}

func (api API) Compressor(compressor string) Option {
	return func(req *readAt.Request) error {
		req.Compressor = compressor
		return nil
	}
}

func (api API) WithResolver(resolver xhttp.ResponseResolver) Option {
	return func(req *readAt.Request) error {
		req.Resolver = xhttp2.ResponseResolver(resolver)
		return nil
	}
}

func defaultResolver(resp *http.Response, ret *readAt.Response) error {
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		body, _ := io.ReadAll(resp.Body)
		return errors.Errorf("http: %v, body: %v", resp.Status, string(body))
	}
	var err error
	ret.Data, err = io.ReadAll(resp.Body)
	if err != nil {
		return errors.Errorf("get data error: %v", err)
	}

	return err
}
