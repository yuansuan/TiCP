package close

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"
	"net/http"
)

type API func(options ...Option) (*session.ApiCloseResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*session.ApiCloseResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}

		uri := "/api/sessions/"
		if req.SessionId != nil {
			uri += *req.SessionId
		}
		uri += "/close"

		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI(uri))

		ret := new(session.ApiCloseResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*session.ApiCloseRequest, error) {
	req := &session.ApiCloseRequest{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}
