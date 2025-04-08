package impl

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/command"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/metrics"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/hpc"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/service/scheduler"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/util"
)

type nodeServiceImpl struct {
	nodeDao dao.NodeDao
	openapi *openapi.OpenAPI
	parser  scheduler.Parser
}

func NewNodeService() (service.NodeService, error) {
	api, err := openapi.NewLocalHPCAPI()
	if err != nil {
		return nil, err
	}

	schedulerType := config.GetConfig().Scheduler.Type
	parser := scheduler.CreateParser(schedulerType)
	nodeService := &nodeServiceImpl{
		nodeDao: dao.NewNodeDao(),
		openapi: api,
		parser:  parser,
	}
	return nodeService, nil
}

func (s *nodeServiceImpl) GetNodeInfo(ctx context.Context, nodeName string) (*dto.NodeDetail, error) {
	logging.Default().Infof("GetNodeInfo method start, nodeName: %s", nodeName)

	result, err := sendPrometheus(ctx, nodeName)
	if err != nil {
		logging.Default().Errorf("sendPrometheusSer method error, err: %v", err)
		return nil, err
	}

	nodeDetail := dto.NodeDetail{}

	buildData(&nodeDetail, nodeName, result)

	return &nodeDetail, nil
}

func GetPrometheusData(ctx context.Context, targetKeys []string) (map[string]model.Value, error) {
	logger := logging.GetLogger(ctx)

	results := make(map[string]model.Value)
	for _, targetKey := range targetKeys {
		queryCondition := fmt.Sprintf("%s", targetKey)
		result, err := util.PromQuery(queryCondition)
		if err != nil {
			logger.Errorf("Error querying Prometheus for target '%s': %v", targetKey, err)
		}

		if result.String() == "" {
			logger.Warnf("the prometheus '%v' target data is empty", targetKey)
		}
		results[targetKey] = result
	}

	return results, nil
}

func (s *nodeServiceImpl) GetNodeList(ctx context.Context, nodeName string, index, size int64) ([]*dto.NodeInfo, int64, error) {
	logging.Default().Infof("GetNodeList method start, nodeName: %s", nodeName)

	nodes, total, err := s.nodeDao.GetNodes(ctx, nodeName, index, size)
	if err != nil {
		logging.Default().Errorf("GetNodeList method error, err: %v", err)
		return nil, 0, err
	}

	dtoList := make([]*dto.NodeInfo, 0, len(nodes))
	for _, node := range nodes {
		nodeDto := dto.NodeInfo{}
		nodeDto.Id = node.Id.String()
		nodeDto.NodeName = node.NodeName
		nodeDto.SchedulerStatus = node.SchedulerStatus
		nodeDto.Status = node.Status
		nodeDto.NodeType = node.NodeType
		nodeDto.QueueName = node.QueueName

		nodeDto.TotalCoreNum = node.TotalCoreNum
		nodeDto.UsedCoreNum = node.UsedCoreNum
		nodeDto.FreeCoreNum = node.FreeCoreNum

		nodeDto.TotalMem = node.TotalMem
		nodeDto.UsedMem = node.UsedMem
		nodeDto.FreeMem = node.FreeMem
		nodeDto.AvailableMem = node.AvailableMem

		nodeDto.CreateTime = node.CreateTime
		dtoList = append(dtoList, &nodeDto)
	}

	return dtoList, total, nil
}

func (s *nodeServiceImpl) NodeOperate(ctx context.Context, nodeNames []string, operation string) error {
	logger := logging.Default()

	//获取节点信息
	nodes, err := s.nodeDao.GetNodeByNames(ctx, nodeNames)
	if err != nil {
		logger.Errorf("get node info failed ，errMsg:[%v]", err)
		return err
	}

	for _, node := range nodes {
		if checkStatus(node.Status, operation) {
			logger.Infof("[%s] status is [%s], operation is [%s], no need to operate", node.NodeName, node.Status, operation)
			continue
		}

		//获取对应命令
		cmd := s.parser.GetCommand(node.NodeName, operation)
		if cmd == "" {
			logger.Errorf("get command failed ，operation:[%s]", operation)
			return status.Error(errcode.ErrOperateTypeNotExist, errcode.MsgOperateTypeNotExist)
		}

		//调用paas平台命令获取节点信息
		customConfig := config.GetConfig()
		resp, err := hpc.Command(ctx, s.openapi, &command.SystemPostRequest{
			Command: cmd,
			Timeout: customConfig.TimeOut, //超时时间
		})
		if err != nil {
			logger.Errorf("invoke openapi.Command failed ，errMsg:[%v]", err)
			return err
		}
		if resp == nil || resp.Data == nil || resp.Data.IsTimeout || resp.Data.Stderr != "" {
			logger.Errorf("invoke openapi.Command failed，Stderr:[%v]", resp.Data.Stderr)
			return errors.New(resp.Data.Stderr)
		}
	}

	return nil
}

func (s *nodeServiceImpl) NodeCoreNum(ctx context.Context) (*dto.CoreStatistics, error) {
	logger := logging.Default()
	nodeList, err := s.nodeDao.NodeList(ctx)
	if err != nil {
		logger.Errorf("get node list failed ，errMsg:[%v]", err)
		return nil, err
	}
	var totalCore, freeCore int
	for _, nodeInfo := range nodeList {
		totalCore += nodeInfo.TotalCoreNum
		if isAvailable(nodeInfo.Status) {
			freeCore += nodeInfo.FreeCoreNum
		}
	}
	response := dto.CoreStatistics{
		TotalNum: totalCore,
		FreeNum:  freeCore,
	}
	return &response, nil
}

func checkStatus(nowStatus, operation string) bool {
	switch operation {
	case consts.NodeClose:
		if nowStatus == consts.Downed || nowStatus == consts.OffLine {
			return true
		} else {
			return false
		}
	case consts.NodeStart:
		if nowStatus != consts.Downed && nowStatus != consts.OffLine {
			return true
		} else {
			return false
		}
	default:
		logging.Default().Errorf("No operation of type [%s] exists", operation)
	}
	return false
}

func buildData(in *dto.NodeDetail, nodeName string, result model.Value) {

	mapValue, err := getMapValue(result)
	if err != nil {
		logging.Default().Errorf("node metric values is empty error, err: %v", err)
		return
	}
	in.NodeName = nodeName                              //节点名称
	in.MaxMem = mapValue[metrics.MemoryTotal]           //最大内存
	in.AvailableMem = mapValue[metrics.MemoryAvailable] //可用内存
	in.UsedMem = mapValue[metrics.MemoryUsed]           //已用内存
	in.FreeMem = mapValue[metrics.MemoryFree]           //空闲内存
	in.FreeSwap = mapValue[metrics.MemorySwapFree]      //空闲交换空间
	in.MaxSwap = mapValue[metrics.MemorySwapTotal]      //最大交换空间
	in.FreeTmp = mapValue[metrics.MemoryTmpFree]        //空闲临时空间
	in.MaxTmp = mapValue[metrics.MemoryTmpTotal]        //最大临时空间

	in.NCore = int64(mapValue[metrics.CPUCore])                                           //核心数量
	in.CpuPercent = mapValue[metrics.CPUPercent]                                          //CPU使用率
	in.UsedCore = int64(mapValue[metrics.CPUCore]) - int64(mapValue[metrics.CPUIdleCore]) //CPU使用核数
	in.CPUIdleCore = int64(mapValue[metrics.CPUIdleCore])                                 //空闲核数
	in.CpuIdleTime = mapValue[metrics.CPUIdleTime]                                        //空闲时间

	in.DiskWriteThroughput = mapValue[metrics.DiskWriteThroughput] //磁盘吞吐率
	in.DiskReadThroughput = mapValue[metrics.DiskReadThroughput]   //磁盘吞吐率

	in.R15m = mapValue[metrics.Load15m] //15分钟负载
	in.R5m = mapValue[metrics.Load5m]   //5分钟负载
	in.R1m = mapValue[metrics.Load1m]   //1分钟负载

}

func strToInt(str string) int64 {
	if str == "" {
		return 0
	}
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		logging.Default().Errorf("convert fail:%v", err)
		return 0
	}
	return num
}

func getMapValue(input model.Value) (map[string]float64, error) {

	result := make(map[string]float64)

	metricValues := input.(model.Vector)
	if len(metricValues) == 0 {
		return nil, errors.New("node metric values is empty")
	}

	for _, v := range metricValues {
		name := string(v.Metric["name"])
		result[name] = float64(v.Value)
	}
	return result, nil
}

func sendPrometheus(ctx context.Context, nodeName string) (model.Value, error) {
	logger := logging.GetLogger(ctx)
	// 待查询的指标
	queryCondition := fmt.Sprintf("%s{host_name=\"%s\"} or %s{host_name=\"%s\"} or %s{host_name=\"%s\"} or %s{host_name=\"%s\"}",
		metrics.MemoryMetrics, nodeName, metrics.CPUMetrics, nodeName,
		metrics.DiskMetrics, nodeName, metrics.AvgLoadMetrics, nodeName)
	result, err := util.PromQuery(queryCondition)
	if err != nil {
		logger.Errorf("Error querying Prometheus for nodeName '%s': %v", nodeName, err)
	}

	if result.String() == "" {
		logger.Warnf("the node '%v' monitor data is empty", nodeName)
	}

	return result, nil
}
