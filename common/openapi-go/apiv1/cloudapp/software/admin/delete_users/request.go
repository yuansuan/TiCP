package delete_users

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/software"
)

type API func(options ...Option) (*software.AdminDeleteUsersResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*software.AdminDeleteUsersResponse, error) {
		req := NewRequest(options)

		resolver := hc.Prepare(xhttp.NewDeleteRequestBuilder().
			URI("/admin/softwares/users").
			Json(req))

		res := new(software.AdminDeleteUsersResponse)
		err := resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, res)
		})

		return res, err
	}
}

func NewRequest(options []Option) *software.AdminDeleteUsersRequest {
	req := new(software.AdminDeleteUsersRequest)

	for _, option := range options {
		option(req)
	}

	return req
}
