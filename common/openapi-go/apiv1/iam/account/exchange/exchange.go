package exchange

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/apiv1/iam/account/add"
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
			URI("/iam/v1/api/account/exchange").
			Json(req))

		ret := new(Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*iam_api.ExchangeCredentialsRequest, error) {
	req := &iam_api.ExchangeCredentialsRequest{}
	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

type Option func(req *iam_api.ExchangeCredentialsRequest) error

func (api API) YsID(ysId string) Option {
	return func(req *iam_api.ExchangeCredentialsRequest) error {
		req.YsId = ysId
		return nil
	}
}

func (api API) Email(email string) Option {
	return func(req *iam_api.ExchangeCredentialsRequest) error {
		req.Email = email
		return nil
	}
}

func (api API) Phone(phone string) Option {
	return func(req *iam_api.ExchangeCredentialsRequest) error {
		req.Phone = phone
		return nil
	}
}

func (api API) Password(password string) Option {
	return func(req *iam_api.ExchangeCredentialsRequest) error {
		req.Password = password
		return nil
	}
}

type Data struct {
	*iam_api.ExchangeCredentialsResponse `json:",inline"`
}

type Response struct {
	add.BaseResponse `json:",inline"`
	Data             *Data `json:"data,omitempty"`
}
