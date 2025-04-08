package poweroff

import (
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(options ...Option) (*session.PowerOffResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*session.PowerOffResponse, error) {
		req := NewRequest(options)

		uri := "/api/sessions/"
		if req.SessionId != nil {
			uri += *req.SessionId
		}
		uri += "/stop"

		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI(uri))
		ret := new(session.PowerOffResponse)
		err := resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) *session.PowerOffRequest {
	req := new(session.PowerOffRequest)

	for _, option := range options {
		option(req)
	}

	return req
}
