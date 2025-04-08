package download

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	xhttp2 "github.com/yuansuan/ticp/common/project-root-api/pkg/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/download"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

type API func(options ...Option) (*download.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*download.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.PrepareDownload(xhttp.NewRequestBuilder().
			URI("/api/storage/download").
			Ref(req))

		ret := new(download.Response)

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

func NewRequest(options []Option) (*download.Request, error) {
	req := &download.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *download.Request) error

func (api API) Path(path string) Option {
	return func(req *download.Request) error {
		req.Path = path
		return nil
	}
}

func (api API) Range(start, end int64) Option {
	return func(req *download.Request) error {
		req.Range = fmt.Sprintf("bytes=%v-%v", start, end)
		return nil
	}
}

func (api API) WithResolver(resolver xhttp.ResponseResolver) Option {
	return func(req *download.Request) error {
		req.Resolver = xhttp2.ResponseResolver(resolver)
		return nil
	}
}

var filenameRegex = regexp.MustCompile(`attachment; filename="(.*?)"`)

func defaultResolver(resp *http.Response, ret *download.Response) error {
	disposition := resp.Header.Get("Content-Disposition")
	arr := filenameRegex.FindStringSubmatch(disposition)
	if len(arr) >= 2 {
		ret.Filename = arr[1]
	}
	ret.FileType = resp.Header.Get("Content-Type")

	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return errors.Errorf("get Content-Length error: %v,Content-Length: %v", err, resp.Header.Get("Content-Length"))
	} else {
		ret.FileSize = int64(size)
	}

	ret.Data, err = io.ReadAll(resp.Body)
	if err != nil {
		return errors.Errorf("get data error: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return errors.Errorf("http: %v, body: %v", resp.Status, string(ret.Data))
	}

	return err
}
