package get_url

import (
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/remoteapp"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(options ...Option) (*remoteapp.ApiGetResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*remoteapp.ApiGetResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}

		uri := "/api/sessions/"
		if req.SessionId != nil {
			uri += *req.SessionId
		}
		uri += "/remoteapps/"
		if req.RemoteAppName != nil {
			uri += *req.RemoteAppName
		}

		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI(uri))

		ret := new(remoteapp.ApiGetResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*remoteapp.ApiGetRequest, error) {
	req := &remoteapp.ApiGetRequest{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}
