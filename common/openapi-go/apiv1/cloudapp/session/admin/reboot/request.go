package reboot

import (
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(options ...Option) (*session.RebootResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*session.RebootResponse, error) {
		req := NewRequest(options)

		uri := "/admin/sessions/"
		if req.SessionId != nil {
			uri += *req.SessionId
		}
		uri += "/restart"

		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI(uri))
		ret := new(session.RebootResponse)
		err := resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) *session.RebootRequest {
	req := new(session.RebootRequest)

	for _, option := range options {
		option(req)
	}

	return req
}
