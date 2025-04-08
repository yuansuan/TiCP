package paymentreduce

import (
	"fmt"
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/paymentreduce"
)

type API func(options ...Option) (*paymentreduce.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*paymentreduce.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			Method(http.MethodPatch).
			URI(fmt.Sprintf("/internal/accounts/:%s/freezedAmount", req.AccountID)).
			Json(req))

		ret := new(paymentreduce.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*paymentreduce.Request, error) {
	req := &paymentreduce.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *paymentreduce.Request) error

func (api API) AccountID(accountID string) Option {
	return func(req *paymentreduce.Request) error {
		req.AccountID = accountID
		return nil
	}
}

func (api API) Comment(comment string) Option {
	return func(req *paymentreduce.Request) error {
		req.Comment = comment
		return nil
	}
}

func (api API) TradeID(tradeID string) Option {
	return func(req *paymentreduce.Request) error {
		req.TradeID = tradeID
		return nil
	}
}

func (api API) MerchandiseID(merchandiseID string) Option {
	return func(req *paymentreduce.Request) error {
		req.MerchandiseID = merchandiseID
		return nil
	}
}

func (api API) MerchandiseName(merchandiseName string) Option {
	return func(req *paymentreduce.Request) error {
		req.MerchandiseName = merchandiseName
		return nil
	}
}

func (api API) UnitPrice(unitPrice int64) Option {
	return func(req *paymentreduce.Request) error {
		req.UnitPrice = unitPrice
		return nil
	}
}

func (api API) PriceDes(priceDes string) Option {
	return func(req *paymentreduce.Request) error {
		req.PriceDes = priceDes
		return nil
	}
}

func (api API) Quantity(quantity float64) Option {
	return func(req *paymentreduce.Request) error {
		req.Quantity = quantity
		return nil
	}
}

func (api API) QuantityUnit(quantityUnit string) Option {
	return func(req *paymentreduce.Request) error {
		req.QuantityUnit = quantityUnit
		return nil
	}
}

func (api API) ResourceID(resourceID string) Option {
	return func(req *paymentreduce.Request) error {
		req.ResourceID = resourceID
		return nil
	}
}

func (api API) StartTime(startTime string) Option {
	return func(req *paymentreduce.Request) error {
		req.StartTime = startTime
		return nil
	}
}

func (api API) EndTime(endTime string) Option {
	return func(req *paymentreduce.Request) error {
		req.EndTime = endTime
		return nil
	}
}

func (api API) Ext(ext string) Option {
	return func(req *paymentreduce.Request) error {
		req.Ext = ext
		return nil
	}
}

func (api API) AccountCashVoucherIDs(accountCashVoucherIDs string) Option {
	return func(req *paymentreduce.Request) error {
		req.AccountCashVoucherIDs = accountCashVoucherIDs
		return nil
	}
}

func (api API) VoucherConsumeMode(voucherConsumeMode int64) Option {
	return func(req *paymentreduce.Request) error {
		req.VoucherConsumeMode = voucherConsumeMode
		return nil
	}
}
