package impl

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	promodel "github.com/prometheus/common/model"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/command"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/metrics"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/hpc"
	pbjob "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/job"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/service/scheduler"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/util"
)

const (
	GetDiskInfoCommand = "df --block-size=1G %s"
)

type dashBoardServiceImpl struct {
	nodeDao   dao.NodeDao
	ReportDao dao.ReportDao
	openapi   *openapi.OpenAPI
	parser    scheduler.Parser
	rpc       *client.GRPC
}

func NewDashBoardService() (service.DashboardService, error) {
	api, err := openapi.NewLocalHPCAPI()
	if err != nil {
		return nil, err
	}

	schedulerType := config.GetConfig().Scheduler.Type
	parser := scheduler.CreateParser(schedulerType)
	rpc := client.GetInstance()

	dashBoardService := &dashBoardServiceImpl{
		nodeDao:   dao.NewNodeDao(),
		openapi:   api,
		parser:    parser,
		ReportDao: dao.NewReportDao(),
		rpc:       rpc,
	}
	return dashBoardService, nil
}

func (s *dashBoardServiceImpl) GetClusterInfo(ctx context.Context) (*dto.ClusterInfo, []*dto.NodeDetail, *dto.Disk, error) {
	//1.获取集群信息
	clusterName, err := s.parser.GetClusterName(ctx, s.openapi)
	if err != nil {
		return nil, nil, nil, err
	}

	//获取所有节点列表
	nodeList, err := s.nodeDao.NodeList(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	//过滤掉配置文件中的节点（计算正常节点和总结点）
	customConfig := config.GetConfig()
	newNodeList, nodeNames := removeHiddenNode(nodeList, customConfig.HiddenNode)
	var totalCore, usedCore, freeCore, availableNodeNum int
	for _, nodeInfo := range newNodeList {
		totalCore += nodeInfo.TotalCoreNum
		usedCore += nodeInfo.UsedCoreNum
		if isAvailable(nodeInfo.Status) {
			freeCore += nodeInfo.FreeCoreNum
			availableNodeNum++
		}
	}

	//2.获取节点信息
	nodeMetricMap, err := getNodeMetricMap(nodeNames)
	if err != nil {
		return nil, nil, nil, err
	}
	nodeDetails := buildNodeData(newNodeList, nodeMetricMap)

	//3.获取磁盘信息
	fields, diskMaps, err := GetDiskData(ctx, s.openapi)
	if err != nil {
		return nil, nil, nil, err
	}

	return &dto.ClusterInfo{
			ClusterName:      strings.TrimSpace(clusterName),
			Cores:            totalCore,
			UsedCores:        usedCore,
			FreeCores:        freeCore,
			AvailableNodeNum: availableNodeNum,
			TotalNodeNum:     len(nodeNames),
		},
		nodeDetails,
		&dto.Disk{
			Fields: fields,
			Data:   diskMaps,
		},
		nil

}

func GetDiskData(ctx context.Context, openapi *openapi.OpenAPI) ([]string, []map[string]interface{}, error) {
	//构建命令，获取配置文件中的挂载路径
	mountPath := config.GetConfig().Scheduler.MountPath
	if mountPath == "" {
		return nil, nil, errors.New("mount path is empty")
	}
	join := strings.Join(strings.Split(mountPath, ","), " ")
	common := fmt.Sprintf(GetDiskInfoCommand, join)

	//获取命令执行结果
	diskData, err := execCommand(ctx, openapi, common)
	if err != nil {
		return nil, nil, err
	}

	//封装数据并返回
	fields, diskMaps := dealWithDiskData(diskData)
	return fields, diskMaps, nil
}

func dealWithDiskData(data string) (fields []string, diskMaps []map[string]interface{}) {
	//获取每一行数据
	rows := strings.Split(data, "\n")
	for _, row := range rows {
		usedMap := make(map[string]interface{})
		freeMap := make(map[string]interface{})
		//跳过第一行
		if strings.Contains(row, "文件系统") {
			continue
		}

		columns := strings.Fields(row)
		//获取每一列的数据
		if len(columns) == 6 {
			used := columns[2]
			free := columns[3]
			path := columns[5]
			fields = append(fields, path)
			usedMap[path] = strToInt(used)
			usedMap[consts.Name] = consts.Used
			freeMap[path] = strToInt(free)
			freeMap[consts.Name] = consts.UnUsed
			diskMaps = append(diskMaps, usedMap, freeMap)
		}
	}
	return fields, diskMaps
}

func execCommand(ctx context.Context, openapi *openapi.OpenAPI, common string) (string, error) {
	logger := logging.Default()

	customConfig := config.GetConfig()

	//获取命令
	resp, err := hpc.Command(ctx, openapi, &command.SystemPostRequest{
		Command: common,
		Timeout: customConfig.TimeOut, //超时时间
	})
	if err != nil {
		logger.Errorf("invoke openapi.Command failed ，errMsg:[%v]", err)
		return "", status.Error(errcode.ErrInternalServer, errcode.MsgInternalServer)
	}
	if resp == nil || resp.Data == nil || resp.Data.IsTimeout || resp.Data.Stderr != "" {
		logger.Errorf("invoke openapi.Command failed , resp:[%v] ", resp)
		return "", status.Error(errcode.ErrInternalServer, errcode.MsgInternalServer)
	}

	return resp.Data.Stdout, nil
}

func buildNodeData(nodeList []*model.NodeInfo, mapValue map[string]*NodeMetric) []*dto.NodeDetail {
	nodeDetails := make([]*dto.NodeDetail, 0, len(nodeList))
	for _, node := range nodeList {
		nodeName := node.NodeName
		nodeDetail := dto.NodeDetail{}
		nodeDetail.NodeName = nodeName                                      //节点名称
		nodeDetail.SchedulerStatus = node.SchedulerStatus                   //调度状态
		nodeDetail.Status = node.Status                                     //调度状态
		nodeDetail.NodeStatus = convertNodeStatus(isAvailable(node.Status)) //节点状态

		nodeDetail.MaxMem = getValue(nodeName, metrics.MemoryTotal, mapValue)           //最大内存
		nodeDetail.AvailableMem = getValue(nodeName, metrics.MemoryAvailable, mapValue) //可用内存
		nodeDetail.UsedMem = getValue(nodeName, metrics.MemoryUsed, mapValue)           //已用内存
		nodeDetail.FreeMem = getValue(nodeName, metrics.MemoryFree, mapValue)           //空闲内存
		nodeDetail.FreeSwap = getValue(nodeName, metrics.MemorySwapFree, mapValue)      //空闲交换空间
		nodeDetail.MaxSwap = getValue(nodeName, metrics.MemorySwapTotal, mapValue)      //最大交换空间
		nodeDetail.FreeTmp = getValue(nodeName, metrics.MemoryTmpFree, mapValue)        //空闲临时空间
		nodeDetail.MaxTmp = getValue(nodeName, metrics.MemoryTmpTotal, mapValue)        //最大临时空间

		nodeDetail.NCore = int64(getValue(nodeName, metrics.CPUCore, mapValue))                                                               //核心数量
		nodeDetail.CpuPercent = getValue(nodeName, metrics.CPUPercent, mapValue)                                                              //CPU使用率
		nodeDetail.UsedCore = int64(getValue(nodeName, metrics.CPUCore, mapValue)) - int64(getValue(nodeName, metrics.CPUIdleCore, mapValue)) //CPU使用时间
		nodeDetail.CPUIdleCore = int64(getValue(nodeName, metrics.CPUIdleCore, mapValue))                                                     //空闲核数
		nodeDetail.CpuIdleTime = getValue(nodeName, metrics.CPUIdleTime, mapValue)                                                            //空闲时间

		nodeDetail.DiskWriteThroughput = getValue(nodeName, metrics.DiskWriteThroughput, mapValue) //磁盘吞吐率
		nodeDetail.DiskReadThroughput = getValue(nodeName, metrics.DiskReadThroughput, mapValue)   //磁盘吞吐率

		nodeDetail.R15m = getValue(nodeName, metrics.Load15m, mapValue) //15分钟负载
		nodeDetail.R5m = getValue(nodeName, metrics.Load5m, mapValue)   //5分钟负载
		nodeDetail.R1m = getValue(nodeName, metrics.Load1m, mapValue)   //1分钟负载
		nodeDetails = append(nodeDetails, &nodeDetail)
	}
	return nodeDetails
}

func convertNodeStatus(status bool) string {
	if status {
		return metrics.Up
	} else {
		return metrics.Down
	}
}

func getValue(nodeName, metric string, mapValue map[string]*NodeMetric) float64 {
	customConfig := config.GetConfig()
	if customConfig.HostNameMapping.Enable {
		//替换映射文件中的名字
		hostNameMap := util.GetHostNameMap()
		if value, ok := hostNameMap[nodeName]; ok {
			nodeName = strings.TrimSpace(value)
		}
	}
	if mapValue[getNodeMetricMapKey(nodeName, metric)] == nil {
		return 0
	}

	return mapValue[getNodeMetricMapKey(nodeName, metric)].Value
}

func isAvailable(status string) bool {
	if consts.Downed == status {
		return false
	}
	return true
}

func removeHiddenNode(list []*model.NodeInfo, hiddenNodeStr string) (nodeList []*model.NodeInfo, nodeNames []string) {
	for _, node := range list {
		if !strings.Contains(hiddenNodeStr, node.NodeName) {
			nodeNames = append(nodeNames, node.NodeName)
			nodeList = append(nodeList, node)
		}
	}
	return
}
func (s *dashBoardServiceImpl) GetResourceInfo(ctx context.Context, req *dto.Request) ([]*dto.ValueStruct, []*dto.ValueStruct, []*dto.ValueStruct, error) {
	timeRange := &dto.TimeRange{
		StartTime: req.Start,
		EndTime:   req.End,
		TimeStep:  consts.RangeTenMinSec,
	}
	//获取CPU利用率
	logger := logging.GetLogger(ctx)
	cpuUtAvg, err := s.ReportDao.GetHostResourceMetricAvgUT(ctx, consts.CPUUtAvg, "", timeRange)
	if err != nil {
		logger.Errorf("err: %v", err)
		return nil, nil, nil, err
	}

	//获取内存利用率
	memUtAvg, err := s.ReportDao.GetHostResourceMetricAvgUT(ctx, consts.MemUtAvg, "", timeRange)
	if err != nil {
		logger.Errorf("err: %v", err)
		return nil, nil, nil, err
	}

	//获取总磁盘IO速率
	totalIoUtAvg, err := s.ReportDao.GetHostResourceMetricAvgUT(ctx, consts.TotalIoUtAvg, "", timeRange)
	if err != nil {
		logger.Errorf("err: %v", err)
		return nil, nil, nil, err
	}
	//获取读磁盘IO速率
	readIoUtAvg, err := s.ReportDao.GetHostResourceMetricAvgUT(ctx, consts.ReadIoUtAvg, "", timeRange)
	if err != nil {
		logger.Errorf("err: %v", err)
		return nil, nil, nil, err
	}
	//获取写磁盘IO速率
	writeIoUtAvg, err := s.ReportDao.GetHostResourceMetricAvgUT(ctx, consts.WriteIoUtAvg, "", timeRange)
	if err != nil {
		logger.Errorf("err: %v", err)
		return nil, nil, nil, err
	}

	//封装数据并返回
	return []*dto.ValueStruct{
			{
				Values: cpuUtAvg,
				Name:   consts.CPUUtAvgName,
			},
		},
		[]*dto.ValueStruct{
			{
				Values: totalIoUtAvg,
				Name:   consts.TotalIoUtAvgName,
			},
			{
				Values: readIoUtAvg,
				Name:   consts.ReadIoUtAvgName,
			},
			{
				Values: writeIoUtAvg,
				Name:   consts.WriteIoUtAvgName,
			},
		},
		[]*dto.ValueStruct{
			{
				Values: memUtAvg,
				Name:   consts.MemUtAvgName,
			},
		}, nil

}
func (s *dashBoardServiceImpl) GetJobInfo(ctx context.Context, req *dto.Request) (jobResRange []*dto.JobStatusValue, jobResLatest []*dto.JobStatusValue, err error) {
	//1.从Prometheus 中获取作业状态
	jobResRange, err = GetJobStatusValue(ctx, &dto.TimeRange{
		StartTime: req.Start,
		EndTime:   req.End,
		TimeStep:  consts.RangeFiveMinSec,
	})
	if err != nil {
		return nil, nil, err
	}
	//2.从Job服务获取作业状态
	resp, err := s.rpc.Job.GetJobStatus(ctx, &pbjob.GetJobStatusRequest{})
	if err != nil {
		return nil, nil, err
	}
	jobResLatest = buildNowJobStatus(resp.JobStatusMap)

	return jobResRange, jobResLatest, nil
}

func buildNowJobStatus(statusMap map[string]int64) []*dto.JobStatusValue {
	var jobStatusValues []*dto.JobStatusValue
	for key, value := range statusMap {
		jobStatusValues = append(jobStatusValues, &dto.JobStatusValue{
			Status:   key,
			JobCount: value,
		})
	}
	return jobStatusValues
}

func GetJobStatusValue(ctx context.Context, timeRange *dto.TimeRange) ([]*dto.JobStatusValue, error) {
	logger := logging.GetLogger(ctx)

	jobStatusValues := make([]*dto.JobStatusValue, 0)
	queryString := fmt.Sprintf("max_over_time(%s{name=~\"^(%s|%s|%s|%s)$\"}[%vs])", consts.JobStatusNum, consts.JobStateCompleted, consts.JobStateFailed, consts.JobStateRunning, consts.JobStatePending, timeRange.TimeStep)
	resp, err := util.PromRangeQuery(queryString, timeRange.StartTime/1000, timeRange.EndTime/1000, timeRange.TimeStep)
	if err != nil {
		return nil, err
	}

	resValues := resp.(promodel.Matrix)

	if len(resValues) <= 0 {
		logger.Warnf("The result of the query for prometheus is empty")
		return jobStatusValues, nil
	}
	for _, promValue := range resValues {
		if promValue.Values != nil && len(promValue.Values) > 0 {
			label := promValue.Metric[consts.Name]
			for _, value := range promValue.Values {
				jobValue := dto.JobStatusValue{}
				jobValue.Timestamp = value.Timestamp.Unix() * 1000
				jobValue.JobCount = int64(value.Value)
				jobValue.Status = string(label)
				jobStatusValues = append(jobStatusValues, &jobValue)
			}
		}
	}

	// 按照时间进行升序排序
	if len(jobStatusValues) >= 2 {
		sort.Slice(jobStatusValues, func(i, j int) bool {
			return jobStatusValues[i].Timestamp < jobStatusValues[j].Timestamp
		})
	}

	return jobStatusValues, nil
}

type NodeMetric struct {
	Name  string
	Value float64
}

func getNodeMetricMap(nodeNames []string) (map[string]*NodeMetric, error) {
	results, err := getPrometheusNodeMap(context.Background(), nodeNames)
	if err != nil {
		return nil, err
	}

	nodeMap := make(map[string]*NodeMetric)
	for _, target := range results {
		metricValues := target.(promodel.Vector)

		for _, v := range metricValues {
			name := string(v.Metric[consts.NameLabel])
			hostName := string(v.Metric[consts.HostNameLabel])
			nodeMetric := &NodeMetric{
				Name:  name,
				Value: float64(v.Value),
			}

			nodeMetricMapKey := getNodeMetricMapKey(hostName, name)
			nodeMap[nodeMetricMapKey] = nodeMetric
		}
	}

	return nodeMap, nil
}
func getNodeMetricMapKey(hostName, name string) string {
	return fmt.Sprintf("%v_%v", hostName, name)
}
func getPrometheusNodeMap(ctx context.Context, nodeNames []string) (map[string]promodel.Value, error) {
	logger := logging.GetLogger(ctx)
	results := make(map[string]promodel.Value)
	enable := config.GetConfig().HostNameMapping.Enable
	for _, nodeName := range nodeNames {
		if enable {
			//替换映射文件中的名字
			hostNameMap := util.GetHostNameMap()
			logger.Infof("hostNameMap '%v'", hostNameMap)
			if value, ok := hostNameMap[nodeName]; ok {
				nodeName = strings.TrimSpace(value)
			}
			logger.Infof("nodeName '%v'", nodeName)
		}
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
		results[nodeName] = result
	}

	return results, nil
}
