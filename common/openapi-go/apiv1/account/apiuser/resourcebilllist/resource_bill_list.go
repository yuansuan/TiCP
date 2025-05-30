package resourcebilllist

import (
	"fmt"
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	resourceBillList "github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/apiuser/resourcebilllist"
)

type API func(options ...Option) (*resourceBillList.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*resourceBillList.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/api/accounts/resourcebilling").
			AddQuery("StartTime", fmt.Sprintf("%v", req.StartTime)).
			AddQuery("EndTime", fmt.Sprintf("%v", req.EndTime)).
			AddQuery("TradeType", fmt.Sprintf("%v", req.TradeType)).
			AddQuery("SignType", fmt.Sprintf("%v", req.SignType)).
			AddQuery("ProductName", fmt.Sprintf("%v", req.ProductName)).
			AddQuery("SortByAsc", fmt.Sprintf("%v", req.SortByAsc)).
			AddQuery("PageSize", fmt.Sprintf("%v", req.PageSize)).
			AddQuery("PageIndex", fmt.Sprintf("%v", req.PageIndex)))

		ret := new(resourceBillList.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*resourceBillList.Request, error) {
	req := &resourceBillList.Request{
		PageIndex: 1,
		PageSize:  10,
	}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *resourceBillList.Request) error

func (api API) PageIndex(index int64) Option {
	return func(req *resourceBillList.Request) error {
		req.PageIndex = index
		return nil
	}
}
func (api API) PageSize(size int64) Option {
	return func(req *resourceBillList.Request) error {
		req.PageSize = size
		return nil
	}
}
func (api API) TradeType(tradeType int64) Option {
	return func(req *resourceBillList.Request) error {
		req.TradeType = tradeType
		return nil
	}
}
func (api API) SignType(signType int64) Option {
	return func(req *resourceBillList.Request) error {
		req.SignType = signType
		return nil
	}
}
func (api API) ProductName(productName string) Option {
	return func(req *resourceBillList.Request) error {
		req.ProductName = productName
		return nil
	}
}
func (api API) SortByAsc(sortByAsc bool) Option {
	return func(req *resourceBillList.Request) error {
		req.SortByAsc = sortByAsc
		return nil
	}
}
func (api API) StartTime(startTime string) Option {
	return func(req *resourceBillList.Request) error {
		req.StartTime = startTime
		return nil
	}
}
func (api API) EndTime(endTime string) Option {
	return func(req *resourceBillList.Request) error {
		req.EndTime = endTime
		return nil
	}
}
