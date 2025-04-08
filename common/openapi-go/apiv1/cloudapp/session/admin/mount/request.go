package mount

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"
)

type API func(opts ...Option) (*session.MountResponse, error)

func New(hc *xhttp.Client) API {
	return func(opts ...Option) (*session.MountResponse, error) {
		req, err := NewRequest(opts)
		if err != nil {
			return nil, err
		}

		uri := "/admin/sessions/"
		if req.SessionId != nil {
			uri += *req.SessionId
		}
		uri += "/mount"

		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI(uri).
			Json(req))

		ret := new(session.MountResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(opts []Option) (*session.MountRequest, error) {
	req := new(session.MountRequest)

	for _, opt := range opts {
		if err := opt(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}
