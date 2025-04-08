package hpc

import (
	"context"
	"fmt"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/openapi-go/apiv1/hpc/command/system/post"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/command"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
)

func Command(ctx context.Context, api *openapi.OpenAPI, req *command.SystemPostRequest) (*command.SystemPostResponse, error) {

	logger := logging.Default()
	logger.Debugf("invoke openapi.Command req:[%+v]", req)

	options := []post.Option{
		api.Client.HPC.Command.System.Execute.Command(req.Command),
		api.Client.HPC.Command.System.Execute.Timeout(req.Timeout),
	}

	response, err := api.Client.HPC.Command.System.Execute(options...)

	if err != nil {
		logger.Errorf("invoke openapi.Command failed ，req:[%+v]，errMsg:[%v]", req, err)
		return nil, err
	}

	if response != nil {
		logging.Default().Debugf("openapi command request id: [%v], req: [%+v]", response.RequestID, req)
	}

	tracelog.Info(ctx, fmt.Sprintf("invoke openapi.Command, command:[%v], response:[%v]", req.Command, response))

	return response, err
}
