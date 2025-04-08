package add

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/remoteapp"
	"net/http"
)

type API func(options ...Option) (*remoteapp.AdminPostResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*remoteapp.AdminPostResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/admin/remoteapps").
			Json(req))

		ret := new(remoteapp.AdminPostResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*remoteapp.AdminPostRequest, error) {
	req := &remoteapp.AdminPostRequest{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}
