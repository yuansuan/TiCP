package storage

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	apimkdir "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/mkdir"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/mkdir"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

// Mkdir 创建文件夹
func Mkdir(api *openapi.OpenAPI, req mkdir.Request) (*mkdir.Response, error) {
	logger := logging.Default()
	logger.Debugf("invoke openapi.Mkdir req:[%+v]", req)

	options := []apimkdir.Option{
		api.Client.Storage.Mkdir.Path(req.Path),
	}

	response, err := api.Client.Storage.Mkdir(options...)

	if err != nil {
		logger.Errorf("invoke openapi.Mkdir failed req:[%+v]，errMsg:[%v]", req, err)
	}

	if response != nil {
		logger.Debugf("openapi mkdir request id: [%v], req: [%+v]", response.RequestID, req)
	}

	return response, err
}
