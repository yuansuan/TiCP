package list

import (
	"net/http"
	"strconv"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(options ...Option) (*session.AdminListResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*session.AdminListResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}

		rb := xhttp.NewRequestBuilder().URI("/admin/sessions")

		if req.PageOffset != nil {
			rb.AddQuery("PageOffset", strconv.Itoa(*req.PageOffset))
		}
		if req.PageSize != nil {
			rb.AddQuery("PageSize", strconv.Itoa(*req.PageSize))
		}
		if req.Status != nil {
			rb.AddQuery("Status", *req.Status)
		}
		if req.SessionIds != nil {
			rb.AddQuery("SessionIds", *req.SessionIds)
		}
		if req.Zone != nil {
			rb.AddQuery("Zone", *req.Zone)
		}
		if req.UserIds != nil {
			rb.AddQuery("UserIds", *req.UserIds)
		}
		if req.WithDeleted == true {
			rb.AddQuery("WithDeleted", "true")
		}

		resolver := hc.Prepare(rb)

		ret := new(session.AdminListResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*session.AdminListRequest, error) {
	req := &session.AdminListRequest{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}
