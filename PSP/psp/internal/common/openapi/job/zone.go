package job

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/zonelist"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

// ZoneList 获取区域列表
func ZoneList(api *openapi.OpenAPI) (*zonelist.Response, error) {
	response, err := api.Client.Job.ZoneList()

	if response != nil {
		logging.Default().Debugf("openapi zone list request id: [%v]", response.RequestID)
	}

	return response, err
}
