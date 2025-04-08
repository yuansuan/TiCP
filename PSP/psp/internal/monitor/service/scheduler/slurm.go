package scheduler

import (
	"context"
	"fmt"
	"strings"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/command"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/hpc"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/consts"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

const (
	// SlurmCloseCommand 关闭节点提交作业
	SlurmCloseCommand = "sudo %sscontrol update nodename=%s state=down reason=\"admin set node down\""

	// SlurmStartCommand 启动节点提交作业
	SlurmStartCommand = "sudo %sscontrol update nodename=%s state=resume reason=\"admin resume node up\""
)

const (
	SlurmStateInfoCommand      = "sinfo -N -o \"%.6N %.5P %.5a %.10T %.15C\""
	GetSlurmClusterNameCommand = "%sscontrol show config | grep ClusterName | cut -d'=' -f2"
	GetSlurmNodeInfoCommand    = "%sscontrol show node"
)

// SlurmParser Slurm调度器对应的解析器
type SlurmParser struct{}

func (s *SlurmParser) Parse(ctx context.Context, openapi *openapi.OpenAPI) ([]map[string]string, error) {
	logger := logging.GetLogger(ctx)
	customConfig := config.GetConfig()

	// 调用openapi 获取节点数据
	resp, err := hpc.Command(ctx, openapi, &command.SystemPostRequest{
		Command: fmt.Sprintf(GetSlurmNodeInfoCommand, strings.TrimSpace(customConfig.Scheduler.CmdPath)),
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

	tracelog.Info(ctx, fmt.Sprintf("daemon slurm parser, command:[scontrol show node], resp:[%v]", resp))

	// 解析数据
	var sliceMap []map[string]string
	dataGroupByNode := strings.Split(resp.Data.Stdout, "\n\n")
	for _, text := range dataGroupByNode {
		if strings.TrimSpace(text) == "" {
			continue
		}
		resultMap := make(map[string]string)
		lines := strings.Split(strings.TrimSpace(text), "\n")
		for _, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), consts.OS) {
				dealWithOS(line, resultMap)
			}

			keyVals := strings.Split(strings.TrimSpace(line), " ")
			for _, keyVal := range keyVals {
				target := strings.Split(strings.TrimSpace(keyVal), "=")
				put(target, resultMap)
			}
		}
		sliceMap = append(sliceMap, resultMap)
	}
	return sliceMap, nil
}

func (s *SlurmParser) GetCommand(nodeName, operation string) (command string) {
	customConfig := config.GetConfig()
	switch operation {
	case consts.NodeClose:
		command = fmt.Sprintf(SlurmCloseCommand, strings.TrimSpace(customConfig.Scheduler.CmdPath), nodeName)
	case consts.NodeStart:
		command = fmt.Sprintf(SlurmStartCommand, strings.TrimSpace(customConfig.Scheduler.CmdPath), nodeName)
	default:
		logging.Default().Errorf("No operation of type [%s] exists", operation)
	}
	return command
}

func (p *SlurmParser) GetStateInfoMap(ctx context.Context, openapi *openapi.OpenAPI) (map[string]string, error) {
	logger := logging.GetLogger(context.Background())
	customConfig := config.GetConfig()

	// 调用openapi 获取节点数据
	resp, err := hpc.Command(ctx, openapi, &command.SystemPostRequest{
		Command: strings.TrimSpace(customConfig.Scheduler.CmdPath) + SlurmStateInfoCommand,
		Timeout: customConfig.TimeOut, //超时时间
	})
	if err != nil {
		logger.Errorf("invoke openapi.Command failed ，errMsg:[%v]", err)
		return nil, err
	}
	if resp == nil || resp.Data == nil || resp.Data.IsTimeout || resp.Data.Stderr != "" {
		logger.Errorf("invoke openapi.Command failed , IsTimeout:[%v] , Stderr:[%v]", resp.Data.IsTimeout, resp.Data.Stderr)
		return nil, status.Error(errcode.ErrInternalServer, errcode.MsgInternalServer)
	}

	// 解析数据
	lines := strings.Split(resp.Data.Stdout, "\n")
	resultMap := make(map[string]string)
	for _, line := range lines {
		// line 表头格式:  NODELIST PARTI AVAIL STATE CPUS(A/I/O/T)
		// 跳过字段说明行
		if strutil.IsEmpty(line) || strings.Contains(line, "NODELIST") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 5 {
			resultMap[CreatKey(strings.TrimSpace(fields[0]), consts.State)] = strings.TrimSpace(fields[3])
			resultMap[CreatKey(strings.TrimSpace(fields[0]), consts.CPUTot)] = processValue(strings.TrimSpace(fields[4]))[3]
			resultMap[CreatKey(strings.TrimSpace(fields[0]), consts.CPUIdle)] = processValue(strings.TrimSpace(fields[4]))[1]
			resultMap[CreatKey(strings.TrimSpace(fields[0]), consts.CPUAlloc)] = processValue(strings.TrimSpace(fields[4]))[0]
		}
	}
	return resultMap, nil
}

func (s *SlurmParser) GetClusterName(ctx context.Context, openapi *openapi.OpenAPI) (string, error) {
	logger := logging.Default()

	customConfig := config.GetConfig()

	//获取命令
	resp, err := hpc.Command(ctx, openapi, &command.SystemPostRequest{
		Command: fmt.Sprintf(GetSlurmClusterNameCommand, strings.TrimSpace(customConfig.Scheduler.CmdPath)),
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

func CreatKey(node, salt string) string {
	return fmt.Sprintf("%s_%s", node, salt)
}

func processValue(cpus string) []string {
	if strutil.IsEmpty(cpus) {
		return []string{"0", "0", "0", "0"}
	}
	split := strings.Split(cpus, "/")
	if len(split) < 4 {
		return []string{"0", "0", "0", "0"}
	}
	return split
}

func dealWithOS(line string, resultMap map[string]string) {
	strSplit := strings.SplitN(strings.TrimSpace(line), " ", 2)
	if len(strSplit) > 1 {
		split := strings.Split(strSplit[0], "=")
		if len(split) > 1 {
			resultMap[consts.OS] = split[1]
		}
	}
}

func put(value []string, in map[string]string) {
	if len(value) < 2 {
		return
	}
	switch value[0] {
	case consts.NodeName:
		in[value[0]] = value[1]
	case consts.CPUTot:
		in[value[0]] = value[1]
	case consts.CPUAlloc:
		in[value[0]] = value[1]
	case consts.State:
		in[value[0]] = value[1]
	case consts.Partitions:
		in[value[0]] = value[1]
	}
}
