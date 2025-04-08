package list

import (
	"net/http"
	"strconv"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/hardware"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(options ...Option) (*hardware.AdminListResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*hardware.AdminListResponse, error) {
		req, err := newRequest(options...)
		if err != nil {
			return nil, err
		}

		rb := xhttp.NewRequestBuilder().URI("/admin/hardwares")
		if req.UserId != nil {
			rb.AddQuery("UserId", *req.UserId)
		}
		if req.Name != nil {
			rb.AddQuery("Name", *req.Name)
		}
		if req.Cpu != nil {
			rb.AddQuery("Cpu", strconv.Itoa(*req.Cpu))
		}
		if req.Mem != nil {
			rb.AddQuery("Mem", strconv.Itoa(*req.Mem))
		}
		if req.Gpu != nil {
			rb.AddQuery("Gpu", strconv.Itoa(*req.Gpu))
		}
		if req.Zone != nil {
			rb.AddQuery("Zone", *req.Zone)
		}
		if req.PageOffset != nil {
			rb.AddQuery("PageOffset", strconv.Itoa(*req.PageOffset))
		}
		if req.PageSize != nil {
			rb.AddQuery("PageSize", strconv.Itoa(*req.PageSize))
		}

		resolver := hc.Prepare(rb)

		ret := new(hardware.AdminListResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func newRequest(options ...Option) (*hardware.AdminListRequest, error) {
	req := &hardware.AdminListRequest{}
	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}
