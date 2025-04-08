package storage

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	apistat "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/stat"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/stat"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

// Stat 查询文件状态
func Stat(api *openapi.OpenAPI, req stat.Request) (*stat.Data, error) {
	logger := logging.Default()
	logger.Debugf("invoke openapi.Stat req:[%+v]", req)

	options := []apistat.Option{
		api.Client.Storage.Stat.Path(req.Path),
	}
	response, err := api.Client.Storage.Stat(options...)

	if err != nil {
		if response.ErrorCode == "PathNotFound" {
			err = status.Error(errcode.ErrFileNotExist, "")
		} else {
			logger.Errorf("invoke openapi.Stat failed req:[%+v]，errMsg:[%v]", req, err)
		}
	}

	if response != nil {
		logger.Debugf("openapi stat request id: [%v], req: [%+v]", response.RequestID, req)
	}

	return response.Data, err
}
