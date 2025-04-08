package execscript

import (
	"net/http"

	"github.com/yuansuan/ticp/common/go-kit/logging/trace"

	"github.com/yuansuan/ticp/common/project-root-api/rdpgo/v1/execscript"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type Option func(req *execscript.Request) error

func (api API) PrivateIP(privateIP string) Option {
	return func(req *execscript.Request) error {
		req.PrivateIP = privateIP
		return nil
	}
}

func (api API) RequestID(requestID string) Option {
	return func(req *execscript.Request) error {
		req.RequestID = requestID
		return nil
	}
}

func (api API) ScriptRunner(scriptRunner string) Option {
	return func(req *execscript.Request) error {
		req.ScriptRunner = scriptRunner
		return nil
	}
}

func (api API) ScriptContentEncoded(scriptContentEncoded string) Option {
	return func(req *execscript.Request) error {
		req.ScriptContentEncoded = scriptContentEncoded
		return nil
	}
}

func (api API) WaitTillEnd(waitTillEnd bool) Option {
	return func(req *execscript.Request) error {
		req.WaitTillEnd = waitTillEnd
		return nil
	}
}

type API func(opts ...Option) (*execscript.Response, error)

func New(hc *xhttp.Client) API {
	return func(opts ...Option) (*execscript.Response, error) {
		req, err := NewRequest(opts...)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			AddQuery("PrivateIP", req.PrivateIP).
			AddHeader(trace.RequestIdKey, req.RequestID).
			URI("/internal/execScript").
			Json(req))

		ret := new(execscript.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(opts ...Option) (*execscript.Request, error) {
	req := new(execscript.Request)

	for _, opt := range opts {
		if err := opt(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}
