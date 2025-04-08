package add_users

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/software"
)

type API func(options ...Option) (*software.AdminPostUsersResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*software.AdminPostUsersResponse, error) {
		req := NewRequest(options)

		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/admin/softwares/users").
			Json(req))

		res := new(software.AdminPostUsersResponse)
		err := resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, res)
		})

		return res, err
	}
}

func NewRequest(options []Option) *software.AdminPostUsersRequest {
	req := new(software.AdminPostUsersRequest)

	for _, option := range options {
		option(req)
	}

	return req
}
