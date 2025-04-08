package storage

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	apicreate "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/create"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/create"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

// Create 创建文件
func Create(api *openapi.OpenAPI, req create.Request) (*create.Response, error) {
	logger := logging.Default()
	logger.Debugf("invoke openapi.Create req:[%+v]", req)

	options := []apicreate.Option{
		api.Client.Storage.Create.Path(req.Path),
		api.Client.Storage.Create.Size(req.Size),
		api.Client.Storage.Create.Overwrite(req.Overwrite),
	}

	response, err := api.Client.Storage.Create(options...)

	if err != nil {
		logger.Errorf("invoke openapi.Create failed req:[%+v]，errMsg:[%v]", req, err)
	}

	if response != nil {
		logger.Debugf("openapi create request id: [%v], req: [%+v]", response.RequestID, req)
	}

	return response, err
}
