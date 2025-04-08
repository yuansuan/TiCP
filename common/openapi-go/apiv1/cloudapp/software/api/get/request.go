package get

import (
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/software"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(options ...Option) (*software.APIGetResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*software.APIGetResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}

		uri := "/api/softwares/"
		if req.SoftwareId != nil {
			uri += *req.SoftwareId
		}

		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI(uri))

		ret := new(software.APIGetResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*software.APIGetRequest, error) {
	req := &software.APIGetRequest{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}
