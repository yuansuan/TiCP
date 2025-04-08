package get

import (
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/quota/admin"
	"net/http"
)

type API func(options ...Option) (*admin.GetStorageQuotaResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*admin.GetStorageQuotaResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		if req.UserID == "" {
			return nil, errors.Errorf("http: 400 Bad Request, " +
				"body: {\"Data\":null,\"ErrorCode\":\"" + commoncode.InvalidUserID + ",\"ErrorMsg\":\"UserID is required\"}")
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/admin/storage/" + req.UserID + "/storageQuota").
			Json(req))

		ret := new(admin.GetStorageQuotaResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*admin.GetStorageQuotaRequest, error) {
	req := &admin.GetStorageQuotaRequest{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *admin.GetStorageQuotaRequest) error

func (api API) UserID(userID string) Option {
	return func(req *admin.GetStorageQuotaRequest) error {
		req.UserID = userID
		return nil
	}
}
