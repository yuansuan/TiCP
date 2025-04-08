package list

import (
	"net/http"
	"strconv"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/software"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(options ...Option) (*software.AdminListResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*software.AdminListResponse, error) {
		req, err := newRequest(options...)
		if err != nil {
			return nil, err
		}
		rb := xhttp.NewRequestBuilder().URI("/admin/softwares")
		if req.UserId != nil {
			rb.AddQuery("UserId", *req.UserId)
		}
		if req.Name != nil {
			rb.AddQuery("Name", *req.Name)
		}
		if req.Platform != nil {
			rb.AddQuery("Platform", *req.Platform)
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

		ret := new(software.AdminListResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func newRequest(options ...Option) (*software.AdminListRequest, error) {
	req := &software.AdminListRequest{}
	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}
