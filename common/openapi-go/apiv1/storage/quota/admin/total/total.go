package total

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/quota/admin"
	"net/http"
)

type API func() (*admin.GetStorageQuotaTotalResponse, error)

func New(hc *xhttp.Client) API {
	return func() (*admin.GetStorageQuotaTotalResponse, error) {

		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/admin/storage/storageQuota/total"))

		ret := new(admin.GetStorageQuotaTotalResponse)
		err := resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}
