package poweron

import (
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(options ...Option) (*session.PowerOnResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*session.PowerOnResponse, error) {
		req := NewRequest(options)

		uri := "/api/sessions/"
		if req.SessionId != nil {
			uri += *req.SessionId
		}
		uri += "/start"

		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI(uri))
		ret := new(session.PowerOnResponse)
		err := resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) *session.PowerOnRequest {
	req := new(session.PowerOnRequest)

	for _, option := range options {
		option(req)
	}

	return req
}
