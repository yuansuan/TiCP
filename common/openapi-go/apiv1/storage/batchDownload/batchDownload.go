package batchDownload

import (
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	xhttp2 "github.com/yuansuan/ticp/common/project-root-api/pkg/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/batchDownload"
	"io"
	"net/http"
	"strconv"
)

type API func(options ...Option) (*batchDownload.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*batchDownload.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.PrepareDownload(xhttp.NewRequestBuilder().
			URI("/api/storage/batchDownload").
			Json(req))

		ret := new(batchDownload.Response)

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

func NewRequest(options []Option) (*batchDownload.Request, error) {
	req := &batchDownload.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *batchDownload.Request) error

func (api API) Paths(paths ...string) Option {
	return func(req *batchDownload.Request) error {
		req.Paths = paths
		return nil
	}
}

func (api API) WithResolver(resolver xhttp.ResponseResolver) Option {
	return func(req *batchDownload.Request) error {
		req.Resolver = xhttp2.ResponseResolver(resolver)
		return nil
	}
}

func (api API) FileName(fileName string) Option {
	return func(req *batchDownload.Request) error {
		req.FileName = fileName
		return nil
	}
}

func (api API) BasePath(basePath string) Option {
	return func(req *batchDownload.Request) error {
		req.BasePath = basePath
		return nil
	}
}

func (api API) IsCompress(isCompress bool) Option {
	return func(req *batchDownload.Request) error {
		req.IsCompress = isCompress
		return nil
	}
}

var zipFileType = "application/zip"

func defaultResolver(resp *http.Response, ret *batchDownload.Response) error {
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		body, _ := io.ReadAll(resp.Body)
		defer func() { _ = resp.Body.Close() }()
		return errors.Errorf("http: %v, body: %v", resp.Status, string(body))
	}

	ret.FileType = resp.Header.Get("Content-Type")

	if resp.Header.Get("Content-Type") != zipFileType {
		size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
		if err != nil {
			return errors.Errorf("get Content-Length error: %v,Content-Length: %v", err, resp.Header.Get("Content-Length"))
		} else {
			ret.FileSize = int64(size)
		}
	}
	return nil
}
