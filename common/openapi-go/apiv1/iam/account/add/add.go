package add

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
)

type API func(options ...Option) (*Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			Method(http.MethodPost).
			URI("/iam/v1/api/account").
			Json(req))

		ret := new(Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			e := common.ParseResp(resp, ret)
			if resp.StatusCode == http.StatusForbidden {
				return nil
			}
			return e
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*iam_api.AddAccountRequest, error) {
	req := &iam_api.AddAccountRequest{}
	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

type Option func(req *iam_api.AddAccountRequest) error

func (api API) Phone(phone string) Option {
	return func(req *iam_api.AddAccountRequest) error {
		req.Phone = phone
		return nil
	}
}

func (api API) Name(name string) Option {
	return func(req *iam_api.AddAccountRequest) error {
		req.Name = name
		return nil
	}
}

func (api API) CompanyName(companyName string) Option {
	return func(req *iam_api.AddAccountRequest) error {
		req.CompanyName = companyName
		return nil
	}
}

func (api API) Password(p string) Option {
	return func(req *iam_api.AddAccountRequest) error {
		req.Password = p
		return nil
	}
}

func (api API) UnifiedSocialCreditCode(code string) Option {
	return func(req *iam_api.AddAccountRequest) error {
		req.UnifiedSocialCreditCode = code
		return nil
	}
}

func (api API) UserChannel(userChannel string) Option {
	return func(req *iam_api.AddAccountRequest) error {
		req.UserChannel = userChannel
		return nil
	}
}

func (api API) Email(email string) Option {
	return func(req *iam_api.AddAccountRequest) error {
		req.Email = email
		return nil
	}
}

type BaseResponse struct {
	RequestId    string `json:"RequestId"`
	ErrorCode    string `json:"ErrorCode"`
	ErrorMessage string `json:"ErrorMessage"`
}

type Data struct {
	*iam_api.AddAccountResponse `json:",inline"`
}

type Response struct {
	BaseResponse `json:",inline"`
	Data         *Data `json:"data,omitempty"`
}
