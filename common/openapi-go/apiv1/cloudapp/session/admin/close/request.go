package close

import (
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(options ...Option) (*session.AdminCloseResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*session.AdminCloseResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}

		uri := "/admin/sessions/"
		if req.SessionId != nil {
			uri += *req.SessionId
		}
		uri += "/close"

		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI(uri).
			Json(req))

		ret := new(session.AdminCloseResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*session.AdminCloseRequest, error) {
	req := &session.AdminCloseRequest{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}
