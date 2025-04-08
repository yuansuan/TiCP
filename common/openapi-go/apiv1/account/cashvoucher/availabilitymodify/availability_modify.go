package availabilitymodify

import (
	"fmt"
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	modify "github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/cashvoucher/availabilitymodify"
)

type API func(options ...Option) (*modify.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*modify.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			Method(http.MethodPatch).
			URI(fmt.Sprintf("/internal/cashvouchers/:%s/availabilitystatus", req.CashVoucherID)).
			Json(req))

		ret := new(modify.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*modify.Request, error) {
	req := &modify.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *modify.Request) error

func (api API) CashVoucherID(cashVoucherID string) Option {
	return func(req *modify.Request) error {
		req.CashVoucherID = cashVoucherID
		return nil
	}
}

func (api API) AvailabilityStatus(availabilityStatus string) Option {
	return func(req *modify.Request) error {
		req.AvailabilityStatus = availabilityStatus
		return nil
	}
}
