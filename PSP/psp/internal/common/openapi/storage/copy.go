package storage

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	apicp "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/copy"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/copy"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func Copy(api *openapi.OpenAPI, req copy.Request) (*copy.Response, error) {

	logger := logging.Default()
	logger.Debugf("invoke openapi.Copy req:[%+v]", req)

	options := []apicp.Option{
		api.Client.Storage.Copy.Src(req.SrcPath),
		api.Client.Storage.Copy.Dest(req.DestPath),
	}

	response, err := api.Client.Storage.Copy(options...)

	if err != nil {
		logger.Errorf("invoke openapi.Copy failed ，req:[%+v]，errMsg:[%v]", req, err)
	}

	if response != nil {
		logger.Debugf("openapi copy request id: [%v], req: [%+v]", response.RequestID, req)
	}

	return response, err

}
