package delete

import (
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(options ...Option) (*session.ApiDeleteResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*session.ApiDeleteResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}

		uri := "/api/sessions/"
		if req.SessionId != nil {
			uri += *req.SessionId
		}

		resolver := hc.Prepare(xhttp.NewDeleteRequestBuilder().
			URI(uri))

		ret := new(session.ApiDeleteResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*session.ApiDeleteRequest, error) {
	req := &session.ApiDeleteRequest{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}
