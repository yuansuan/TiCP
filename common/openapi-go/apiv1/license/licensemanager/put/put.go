package put

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	licmanager "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager"
	"net/http"
)

type API func(options ...Option) (*licmanager.PutLicManagerResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*licmanager.PutLicManagerResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPutRequestBuilder().
			URI("/admin/licenseManagers/" + req.Id).
			Json(req))

		ret := new(licmanager.PutLicManagerResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}
func NewRequest(options []Option) (*licmanager.PutLicManagerRequest, error) {
	req := new(licmanager.PutLicManagerRequest)

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *licmanager.PutLicManagerRequest) error

func (api API) Id(id string) Option {
	return func(req *licmanager.PutLicManagerRequest) error {
		req.Id = id
		return nil
	}
}

func (api API) AppType(appType string) Option {
	return func(req *licmanager.PutLicManagerRequest) error {
		req.AppType = appType
		return nil
	}
}

func (api API) Os(os int) Option {
	return func(req *licmanager.PutLicManagerRequest) error {
		req.Os = int(os)
		return nil
	}
}

func (api API) Desc(desc string) Option {
	return func(req *licmanager.PutLicManagerRequest) error {
		req.Desc = desc
		return nil
	}
}

func (api API) ComputeRule(computeRule string) Option {
	return func(req *licmanager.PutLicManagerRequest) error {
		req.ComputeRule = computeRule
		return nil
	}
}

func (api API) Status(status int) Option {
	return func(req *licmanager.PutLicManagerRequest) error {
		req.Status = status
		return nil
	}
}
