package storage

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	apirealpath "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/realpath"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/realpath"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

// Realpath 查询文件真实路径
func Realpath(api *openapi.OpenAPI, req *realpath.Request) (*realpath.Response, error) {
	logger := logging.Default()
	logger.Debugf("invoke openapi.Realpath req:[%+v]", req)

	options := []apirealpath.Option{
		api.Client.Storage.Realpath.RelativePath(req.RelativePath),
	}

	response, err := api.Client.Storage.Realpath(options...)

	if err != nil {
		logger.Errorf("invoke openapi.Realpath failed ,req:[%+v]，errMsg:[%v]", req, err)
	}

	return response, err
}
