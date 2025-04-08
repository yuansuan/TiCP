package restore

import (
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(option ...Option) (*session.ApiRestoreResponse, error)

func New(hc *xhttp.Client) API {
	return func(option ...Option) (*session.ApiRestoreResponse, error) {
		req := NewRequest(option)

		uri := "/api/sessions/"
		if req.SessionId != nil {
			uri += *req.SessionId
		}
		uri += "/restore"

		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().URI(uri))
		ret := new(session.ApiRestoreResponse)
		err := resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) *session.ApiRestoreRequest {
	req := new(session.ApiRestoreRequest)

	for _, option := range options {
		option(req)
	}

	return req
}
