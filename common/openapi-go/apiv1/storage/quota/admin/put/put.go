package put

import (
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/quota/admin"
	"net/http"
)

type API func(options ...Option) (*admin.PutStorageQuotaResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*admin.PutStorageQuotaResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		if req.UserID == "" {
			return nil, errors.Errorf("http: 400 Bad Request, " +
				"body: {\"Data\":null,\"ErrorCode\":\"" + commoncode.InvalidUserID + ",\"ErrorMsg\":\"UserID is required\"}")
		}
		resolver := hc.Prepare(xhttp.NewPutRequestBuilder().
			URI("/admin/storage/" + req.UserID + "/storageQuota").
			Json(req))

		ret := new(admin.PutStorageQuotaResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*admin.PutStorageQuotaRequest, error) {
	req := &admin.PutStorageQuotaRequest{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *admin.PutStorageQuotaRequest) error

func (api API) UserID(userID string) Option {
	return func(req *admin.PutStorageQuotaRequest) error {
		req.UserID = userID
		return nil
	}
}

func (api API) StorageLimit(storageLimit float64) Option {
	return func(req *admin.PutStorageQuotaRequest) error {
		req.StorageLimit = storageLimit
		return nil
	}
}
