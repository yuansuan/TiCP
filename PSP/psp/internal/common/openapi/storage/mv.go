package storage

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	apimv "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/mv"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/mv"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

// Mv 移动文件
func Mv(api *openapi.OpenAPI, req mv.Request) (*mv.Response, error) {
	logger := logging.Default()
	logger.Debugf("invoke openapi.Mv req:[%+v]", req)

	options := []apimv.Option{
		api.Client.Storage.Mv.Src(req.SrcPath),
		api.Client.Storage.Mv.Dest(req.DestPath),
	}

	response, err := api.Client.Storage.Mv(options...)

	if err != nil {
		logger.Errorf("invoke openapi.Mv failed req:[%+v]，errMsg:[%v]", req, err)
	}

	if response != nil {
		logger.Debugf("openapi mv request id: [%v], req: [%+v]", response.RequestID, req)
	}

	return response, err
}
