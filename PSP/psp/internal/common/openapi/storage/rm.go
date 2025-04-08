package storage

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	apirm "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/rm"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/rm"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

// Rm 删除文件
func Rm(api *openapi.OpenAPI, req rm.Request) (*rm.Response, error) {
	logger := logging.Default()
	logger.Debugf("invoke openapi.Rm req:[%+v]", req)

	options := []apirm.Option{
		api.Client.Storage.Rm.Path(req.Path),
	}

	response, err := api.Client.Storage.Rm(options...)

	if err != nil {
		logger.Errorf("invoke openapi.Rm failed req:[%+v]，errMsg:[%v]", req, err)
	}

	if response != nil {
		logger.Debugf("openapi rm request id: [%v], req: [%+v]", response.RequestID, req)
	}

	return response, err
}
