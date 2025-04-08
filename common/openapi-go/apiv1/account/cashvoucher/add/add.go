package add

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/cashvoucher/add"
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
			URI("/internal/cashvouchers").
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

func (api API) CashVoucherName(cashVoucherName string) Option {
	return func(req *add.Request) error {
		req.CashVoucherName = cashVoucherName
		return nil
	}
}

func (api API) Amount(amount int64) Option {
	return func(req *add.Request) error {
		req.Amount = amount
		return nil
	}
}

func (api API) ExpiredType(expiredType int64) Option {
	return func(req *add.Request) error {
		req.ExpiredType = expiredType
		return nil
	}
}

func (api API) AbsExpiredTime(absExpiredTime string) Option {
	return func(req *add.Request) error {
		req.AbsExpiredTime = absExpiredTime
		return nil
	}
}

func (api API) RelExpiredTime(relExpiredTime int64) Option {
	return func(req *add.Request) error {
		req.RelExpiredTime = relExpiredTime
		return nil
	}
}

func (api API) Comment(comment string) Option {
	return func(req *add.Request) error {
		req.Comment = comment
		return nil
	}
}
