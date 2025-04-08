package create

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/create"
)

type API func(options ...Option) (*create.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*create.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			Method(http.MethodPost).
			URI("/internal/accounts").
			Json(req))
		ret := new(create.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})
		return ret, err
	}
}

func NewRequest(options []Option) (*create.Request, error) {
	req := &create.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

type Option func(req *create.Request) error

func (api API) AccountName(accountName string) Option {
	return func(req *create.Request) error {
		req.AccountName = accountName
		return nil
	}
}

func (api API) UserID(userID string) Option {
	return func(req *create.Request) error {
		req.UserID = userID
		return nil
	}
}

func (api API) AccountType(accountType int64) Option {
	return func(req *create.Request) error {
		req.AccountType = accountType
		return nil
	}
}
