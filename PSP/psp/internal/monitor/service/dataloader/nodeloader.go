package dataloader

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	proModel "github.com/prometheus/common/model"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/metrics"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/service/impl"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/service/scheduler"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
)

// NodeLoader ...
type NodeLoader struct {
	nodeDao dao.NodeDao
	openapi *openapi.OpenAPI
	sid     *snowflake.Node
	parser  scheduler.Parser
}

// NewNodeLoader 节点数据同步
func NewNodeLoader() (*NodeLoader, error) {
	nodeDao := dao.NewNodeDao()

	api, err := openapi.NewLocalHPCAPI()
	if err != nil {
		return nil, err
	}

	node, err := snowflake.GetInstance()
	if err != nil {
		return nil, err
	}

	//根据类型不同new不同的解析器
	schedulerType := config.GetConfig().Scheduler.Type
	parser := scheduler.CreateParser(schedulerType)
	if parser == nil {
		return nil, status.Error(errcode.ErrSchedulerTypeNotSet, errcode.MsgSchedulerTypeNotSet)
	}
	impl := &NodeLoader{
		nodeDao: nodeDao,
		openapi: api,
		sid:     node,
		parser:  parser,
	}

	return impl, nil
}

// NodeLoaderStart 机器节点信息同步
func (loader *NodeLoader) NodeLoaderStart() {
	go loader.NodeGatherTicker()
}

// NodeGatherTicker 机器节点信息同步定时器
func (loader *NodeLoader) NodeGatherTicker() {
	ctx := context.Background()
	logger := logging.GetLogger(ctx)

	// 根据配置启用
	syncData := config.GetConfig().SyncData
	if !syncData.Enable {
		logger.Infof("sync machine node info routine has disabled")
		return
	}

	timerDuration := time.Second * time.Duration(syncData.Interval)
	timer := time.NewTimer(timerDuration)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			if err := loader.NodeGather(); err != nil {
				logger.Errorf("node gather err: %v", err)
			}
			timer.Reset(timerDuration)
		}
	}
}

// NodeGather 机器节点数据采集
func (loader *NodeLoader) NodeGather() error {
	ctx := context.Background()
	logger := logging.GetLogger(ctx)

	//1.构建待插入的节点信息
	nodes, err := loader.buildNodeInfos(ctx)
	if err != nil {
		logger.Errorf("buildNodeInfo failed,err:%v", err)
		return err
	}

	//2.查询所有节点
	nodeList, err := loader.nodeDao.NodeList(ctx)
	if err != nil {
		logger.Errorf("nodeDao.NodeList failed,err:%v", err)
	}

	//3.节点分类
	inserts, updates, deletes := loader.dataSort(nodeList, nodes)
	_, err = boot.MW.DefaultTransaction(ctx, func(session *xorm.Session) (interface{}, error) {
		//新增节点信息
		if len(inserts) > 0 {
			err = loader.nodeDao.AddNodes(ctx, session, inserts)
			if err != nil {
				logger.Errorf("nodeDao.AddNodes failed,err:%v", err)
				return nil, err
			}
		}

		//更新节点信息
		for _, update := range updates {
			err = loader.nodeDao.UpdateNode(ctx, session, update)
			if err != nil {
				logger.Errorf("nodeDao.UpdateNode failed,err:%v", err)
				return nil, err
			}
		}

		//删除节点信息
		if len(deletes) > 0 {
			err = loader.nodeDao.DeleteNotIds(ctx, session, deletes)
			if err != nil {
				logger.Errorf("nodeDao.DeleteNodes failed,err:%v", err)
				return nil, err
			}
		}
		return nil, nil
	})
	if err != nil {
		logger.Errorf("nodeServiceImpl.UpdateNodes failed,err:%v", err)
		return err
	}

	return nil
}

func (loader *NodeLoader) dataSort(nodeList []*model.NodeInfo, nodes []*model.NodeInfo) (inserts []*model.NodeInfo, updates []*model.NodeInfo, deletes []snowflake.ID) {
	nodeMap := make(map[string]snowflake.ID)
	for _, node := range nodeList {
		nodeMap[node.NodeName] = node.Id
	}
	for _, node := range nodes {
		if _, ok := nodeMap[node.NodeName]; !ok {
			node.Id = loader.sid.Generate()
			inserts = append(inserts, node)
			node.CreateTime = time.Now()
			node.UpdateTime = time.Now()
		} else {
			node.Id = nodeMap[node.NodeName]
			updates = append(updates, node)
			node.UpdateTime = time.Now()
		}
		deletes = append(deletes, node.Id)
	}
	return inserts, updates, deletes
}

func (loader *NodeLoader) buildNodeInfos(ctx context.Context) ([]*model.NodeInfo, error) {
	//获取节点信息
	sliceMap, err := loader.parser.Parse(ctx, loader.openapi)
	if err != nil {
		logging.Default().Errorf("parser.Parse failed,err:%v", err)
		return nil, err
	}

	// 待查询的指标
	targetKeys := []string{
		metrics.MemoryMetrics,
		metrics.CPUMetrics,
	}

	nodeMetricMap, err := getNodeMetricMap(targetKeys)
	if err != nil {
		logging.Default().Errorf("get CpuMemInfo failed,err:%v", err)
		return nil, err
	}
	tracelog.Info(ctx, fmt.Sprintf("daemon get prometheus data, data:[%v]", nodeMetricMap))
	//获取状态信息
	nodeInfoMap, err := loader.parser.GetStateInfoMap(ctx, loader.openapi)
	if err != nil {
		logging.Default().Errorf("get SchedRunningInfo failed,err:%v", err)
		return nil, err
	}

	nodes := make([]*model.NodeInfo, 0, len(sliceMap))
	for _, node := range sliceMap {
		totalMem := getNodeMetricValue(metrics.MemoryMetrics, metrics.MemoryTotal, node, nodeMetricMap)
		usedMem := getNodeMetricValue(metrics.MemoryMetrics, metrics.MemoryUsed, node, nodeMetricMap)
		freeMem := getNodeMetricValue(metrics.MemoryMetrics, metrics.MemoryFree, node, nodeMetricMap)
		availableMem := getNodeMetricValue(metrics.MemoryMetrics, metrics.MemoryAvailable, node, nodeMetricMap)
		info := model.NodeInfo{
			NodeName:        node[consts.NodeName],
			NodeType:        node[consts.OS],
			SchedulerStatus: getRightValue(node[consts.State], nodeInfoMap[scheduler.CreatKey(node[consts.NodeName], consts.State)]),
			Status:          getRightValue(statusReflect(node[consts.State]), statusReflect(nodeInfoMap[scheduler.CreatKey(node[consts.NodeName], consts.State)])),
			QueueName:       getQueueValue(node[consts.Queue], node[consts.Partitions]),
			PlatformName:    getRightValue(node[consts.Platform], ""),
			TotalCoreNum:    toInt(getRightValue(node[consts.CPUTot], nodeInfoMap[scheduler.CreatKey(node[consts.NodeName], consts.CPUTot)])),
			UsedCoreNum:     toInt(getRightValue(node[consts.CPUAlloc], nodeInfoMap[scheduler.CreatKey(node[consts.NodeName], consts.CPUAlloc)])),
			TotalMem:        totalMem,
			UsedMem:         usedMem,
			FreeMem:         freeMem,
			AvailableMem:    availableMem,
		}
		info.FreeCoreNum = toInt(getRightValue(strconv.Itoa(info.TotalCoreNum-info.UsedCoreNum), nodeInfoMap[scheduler.CreatKey(node[consts.NodeName], consts.CPUIdle)]))
		nodes = append(nodes, &info)
	}
	return nodes, nil
}

func getQueueValue(pbsProValue, slurmValue string) string {
	defaultQueue := config.GetConfig().Scheduler.DefaultQueue
	schedulerType := config.GetConfig().Scheduler.Type
	switch schedulerType {
	case scheduler.PbsPro:
		if pbsProValue == "" {
			logging.Default().Infof("scheduler node defaultQueue [%s]", defaultQueue)
			return defaultQueue
		}
		return pbsProValue
	case scheduler.Slurm:
		if slurmValue == "" {
			logging.Default().Infof("scheduler node defaultQueue [%s]", defaultQueue)
			return defaultQueue
		}
		return slurmValue
	default:
		logging.Default().Warnf("not support scheduler type [%s]", schedulerType)
	}

	return ""
}

func getRightValue(nodeValue, mapValue string) string {
	schedulerType := config.GetConfig().Scheduler.Type
	switch schedulerType {
	case scheduler.PbsPro:
		return nodeValue
	case scheduler.Slurm:
		return mapValue
	default:
		logging.Default().Warnf("not support scheduler type [%s]", schedulerType)
	}
	return ""
}
func statusReflect(status string) string {
	unavailableStatus := config.GetConfig().UnavailableStatus
	if unavailableStatus == "" || len(unavailableStatus) == 0 {
		//设置默认值
		unavailableStatus = strings.Join([]string{consts.Downed, consts.Drain, consts.OffLine}, ",")
	}
	values := strings.Split(unavailableStatus, ",")
	for _, value := range values {
		if strings.Contains(status, value) {
			return consts.Downed
		}
	}

	return consts.Idle
}

func toInt(str string) int {
	if str == "" {
		return 0
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

func getNodeMetricValue(key, name string, nodeMap map[string]string, nodeMetricMap map[string]*NodeMetric) int {
	customConfig := config.GetConfig()
	nodeName := nodeMap[consts.NodeName]
	if customConfig.HostNameMapping.Enable {
		//替换映射文件中的名字
		hostNameMap := util.GetHostNameMap()
		if value, ok := hostNameMap[nodeName]; ok {
			nodeName = strings.TrimSpace(value)
		}
	}
	nodeMetricMapKey := getNodeMetricMapKey(key, nodeName, name)
	if metricValue, ok := nodeMetricMap[nodeMetricMapKey]; ok {
		return int(metricValue.Value)
	}
	return 0
}

func getNodeMetricMapKey(key, hostName, name string) string {
	return fmt.Sprintf("%v_%v_%v", key, hostName, name)
}

func getNodeMetricMap(targetKeys []string) (map[string]*NodeMetric, error) {
	results, err := impl.GetPrometheusData(context.Background(), targetKeys)
	if err != nil {
		return nil, err
	}

	nodeMap := make(map[string]*NodeMetric)
	for key, target := range results {
		metricValues := target.(proModel.Vector)
		if len(metricValues) == 0 {
			return nil, errors.New("node metric values is empty")
		}

		for _, v := range metricValues {
			name := string(v.Metric["name"])
			hostName := string(v.Metric["host_name"])
			nodeMetric := &NodeMetric{
				Key:       key,
				Name:      name,
				HostName:  hostName,
				Value:     float64(v.Value),
				Timestamp: v.Timestamp.Unix(),
			}

			nodeMetricMapKey := getNodeMetricMapKey(key, hostName, name)
			nodeMap[nodeMetricMapKey] = nodeMetric
		}
	}

	return nodeMap, nil
}

type NodeMetric struct {
	Key       string
	Name      string
	HostName  string
	Value     float64
	Timestamp int64
}
