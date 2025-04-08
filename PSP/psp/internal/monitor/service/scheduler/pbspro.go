package scheduler

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/command"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/hpc"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/consts"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
)

const (
	// PbsProCloseCommand 关闭节点提交作业
	PbsProCloseCommand = "%spbsnodes -o %s"

	// PbsProStartCommand 启动节点提交作业
	PbsProStartCommand = "%spbsnodes -c %s"

	GetPbsProClusterNameCommand = "grep \"PBS_SERVER\" %s | cut -d \"=\" -f 2"
	PbsProNodeInfoCommand       = "%spbsnodes -av"
)

var (
	platform string
	once     = sync.Once{}
)

// PbsProParser Pbs调度器对应的解析器
type PbsProParser struct{}

func (p *PbsProParser) Parse(ctx context.Context, openapi *openapi.OpenAPI) ([]map[string]string, error) {
	logger := logging.GetLogger(ctx)
	customConfig := config.GetConfig()

	// 调用openapi 获取节点信息
	resp, err := hpc.Command(ctx, openapi, &command.SystemPostRequest{
		Command: fmt.Sprintf(PbsProNodeInfoCommand, strings.TrimSpace(customConfig.Scheduler.CmdPath)),
		Timeout: customConfig.TimeOut, //超时时间
	})
	if err != nil {
		logger.Errorf("invoke openapi.Command failed ，errMsg:[%v]", err)
		return nil, err
	}
	if resp == nil || resp.Data == nil || resp.Data.IsTimeout || resp.Data.Stderr != "" {
		logger.Errorf("invoke openapi.Command failed , rsp:[%v]", resp)
		return nil, status.Error(errcode.ErrInternalServer, errcode.MsgInternalServer)
	}

	tracelog.Info(ctx, fmt.Sprintf("daemon pbsPro parser, command:[pbsnodes -av], resp:[%v]", resp))

	// 解析数据
	var sliceMap []map[string]string
	dataGroupByNode := strings.Split(removeAllSpaces(resp.Data.Stdout), "\n\n")
	for _, text := range dataGroupByNode {
		if strings.TrimSpace(text) == "" {
			continue
		}
		resultMap := make(map[string]string)
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			target := strings.Split(line, "=")
			putValue(target, resultMap)
		}
		sliceMap = append(sliceMap, resultMap)
	}

	return sliceMap, nil
}

func (p *PbsProParser) GetStateInfoMap(ctx context.Context, openapi *openapi.OpenAPI) (map[string]string, error) {
	return nil, nil
}
func (p *PbsProParser) GetClusterName(ctx context.Context, openapi *openapi.OpenAPI) (string, error) {
	logger := logging.Default()

	customConfig := config.GetConfig()
	//获取命令
	resp, err := hpc.Command(ctx, openapi, &command.SystemPostRequest{
		Command: fmt.Sprintf(GetPbsProClusterNameCommand, customConfig.Scheduler.ConfigPath),
		Timeout: customConfig.TimeOut, //超时时间
	})
	if err != nil {
		logger.Errorf("invoke openapi.Command failed ，errMsg:[%v]", err)
		return "", status.Error(errcode.ErrInternalServer, errcode.MsgInternalServer)
	}
	if resp == nil || resp.Data == nil || resp.Data.IsTimeout || resp.Data.Stderr != "" {
		logger.Errorf("invoke openapi.Command failed , IsTimeout:[%v] , Stderr:[%v]", resp.Data.IsTimeout, resp.Data.Stderr)
		return "", status.Error(errcode.ErrInternalServer, errcode.MsgInternalServer)
	}
	return resp.Data.Stdout, nil
}

// 去除字符串中的所有空格
func removeAllSpaces(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

func putValue(value []string, in map[string]string) {
	if len(value) < 2 {
		return
	}
	platformName := GetPlatformName()
	switch value[0] {
	case consts.Mom:
		in[consts.NodeName] = value[1]
	case consts.Arch:
		in[consts.OS] = value[1]
	case consts.PbsProState:
		in[consts.State] = value[1]
	case consts.ToTalCpus:
		in[consts.CPUTot] = value[1]
	case consts.Assigned:
		in[consts.CPUAlloc] = value[1]
	case consts.Queue:
		in[consts.Queue] = value[1]
	case platformName:
		in[consts.Platform] = value[1]
	}
}
func (p *PbsProParser) GetCommand(nodeName, operation string) (command string) {
	customConfig := config.GetConfig()
	switch operation {
	case consts.NodeClose:
		command = fmt.Sprintf(PbsProCloseCommand, strings.TrimSpace(customConfig.Scheduler.CmdPath), nodeName)
	case consts.NodeStart:
		command = fmt.Sprintf(PbsProStartCommand, strings.TrimSpace(customConfig.Scheduler.CmdPath), nodeName)
	default:
		logging.Default().Errorf("No operation of type [%s] exists", operation)
	}
	return command
}

func GetPlatformName() string {
	once.Do(func() {
		initPlatformName()
	})
	if platform == "" {
		return consts.DefaultPlatform
	}
	return platform
}

func initPlatformName() {
	platform = config.GetConfig().Scheduler.ResAvailablePlatform
	if platform == "" {
		logging.Default().Warnf("res_available_platform is empty, use default platform name [%s]", consts.DefaultPlatform)
	}
}
