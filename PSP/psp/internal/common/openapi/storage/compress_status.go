package storage

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	compressstatus "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/compress/status"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/compress/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func CompressStatus(api *openapi.OpenAPI, req status.Request) (*status.Data, error) {
	logger := logging.Default()
	logger.Debugf("invoke openapi.compress status req:[%+v]", req)

	options := []compressstatus.Option{
		api.Client.Storage.CompressStatus.CompressID(req.CompressID),
	}
	response, err := api.Client.Storage.CompressStatus(options...)

	if err != nil {
		logger.Errorf("invoke openapi.compress status failed req:[%+v]，errMsg:[%v]", req, err)
	}

	if response != nil {
		logging.Default().Debugf("openapi check compress status id: [%v], req: [%+v]", response.RequestID, req)
	}

	return response.Data, err
}
