package get

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/quota/api"
	"net/http"
)

type API func() (*api.GetStorageQuotaResponse, error)

func New(hc *xhttp.Client) API {
	return func() (*api.GetStorageQuotaResponse, error) {

		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/api/storage/storageQuota"))

		ret := new(api.GetStorageQuotaResponse)
		err := resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}
