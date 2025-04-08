package list

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	licmanager "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager"
	"net/http"
)

type API func(options ...Option) (*licmanager.ListLicManagerResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*licmanager.ListLicManagerResponse, error) {
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/admin/licenseManagers"))

		ret := new(licmanager.ListLicManagerResponse)
		err := resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

type Option func(req interface{}) error
