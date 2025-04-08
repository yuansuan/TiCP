package execscript

import (
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(opts ...Option) (*session.ExecScriptResponse, error)

func New(hc *xhttp.Client) API {
	return func(opts ...Option) (*session.ExecScriptResponse, error) {
		req, err := NewRequest(opts)
		if err != nil {
			return nil, err
		}

		uri := "/api/sessions/"
		if req.SessionId != nil {
			uri += *req.SessionId
		}
		uri += "/execScript"

		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI(uri).
			Json(req))

		ret := new(session.ExecScriptResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(opts []Option) (*session.ExecScriptRequest, error) {
	req := new(session.ExecScriptRequest)
	for _, opt := range opts {
		if err := opt(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}
