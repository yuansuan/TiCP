package add

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/accountcashvoucher/add"
)

type API func(options ...Option) (*add.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*add.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			Method(http.MethodPost).
			URI("/internal/accountcashvouchers").
			Json(req))

		ret := new(add.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*add.Request, error) {
	req := &add.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *add.Request) error

func (api API) CashVoucherID(cashVoucherID string) Option {
	return func(req *add.Request) error {
		req.CashVoucherID = cashVoucherID
		return nil
	}
}

func (api API) AccountIDs(accountIDs string) Option {
	return func(req *add.Request) error {
		req.AccountIDs = accountIDs
		return nil
	}
}
