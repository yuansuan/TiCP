package patch

import (
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/remoteapp"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(options ...Option) (*remoteapp.AdminPatchResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*remoteapp.AdminPatchResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}

		uri := "/admin/remoteapps/"
		if req.RemoteAppId != nil {
			uri += *req.RemoteAppId
		}

		resolver := hc.Prepare(xhttp.NewPatchRequestBuilder().
			URI(uri).
			Json(req))

		ret := new(remoteapp.AdminPatchResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*remoteapp.AdminPatchRequest, error) {
	req := &remoteapp.AdminPatchRequest{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}
