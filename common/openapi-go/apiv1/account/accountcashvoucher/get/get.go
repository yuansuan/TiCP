package get

import (
	"fmt"
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/accountcashvoucher/get"
)

type API func(options ...Option) (*get.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*get.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			Method(http.MethodGet).
			URI(fmt.Sprintf("/internal/accountcashvouchers/:%s", req.AccountCashVoucherID)).
			Json(req))

		ret := new(get.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*get.Request, error) {
	req := &get.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *get.Request) error

func (api API) AccountCashVoucherID(accountCashVoucherID string) Option {
	return func(req *get.Request) error {
		req.AccountCashVoucherID = accountCashVoucherID
		return nil
	}
}
