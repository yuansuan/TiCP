package clean

import (
	"net/http"

	"github.com/yuansuan/ticp/common/go-kit/logging/trace"

	"github.com/yuansuan/ticp/common/project-root-api/rdpgo/v1/clean"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type Option func(req *clean.Request) error

func (api API) PrivateIP(privateIP string) Option {
	return func(req *clean.Request) error {
		req.PrivateIP = privateIP
		return nil
	}
}

func (api API) RequestID(requestID string) Option {
	return func(req *clean.Request) error {
		req.RequestID = requestID
		return nil
	}
}

type API func(opts ...Option) (*clean.Response, error)

func New(hc *xhttp.Client) API {
	return func(opts ...Option) (*clean.Response, error) {
		req, err := NewRequest(opts...)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			AddQuery("PrivateIP", req.PrivateIP).
			AddHeader(trace.RequestIdKey, req.RequestID).
			URI("/internal/clean"))

		ret := new(clean.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(opts ...Option) (*clean.Request, error) {
	req := new(clean.Request)

	for _, opt := range opts {
		if err := opt(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}
