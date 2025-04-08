package rpc

import (
	"context"
	"strings"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/monitor"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/config"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

// QueueList 获取队列列表
func (s *GRPCService) QueueList(ctx context.Context, req *monitor.QueueListRequest) (*monitor.QueueListResponse, error) {
	logger := logging.GetLogger(ctx)
	queueList, err := s.nodeDao.QueueList(ctx)
	if err != nil {
		logger.Errorf("get queue list error, err: %v", err)
		return nil, err
	}

	//过滤队列，排除配置文件中的隐藏队列
	hiddenQueueMap := getHiddenQueue()
	var resultList []string
	if len(hiddenQueueMap) != 0 {
		for _, queue := range queueList {
			if _, ok := hiddenQueueMap[queue]; !ok {
				resultList = append(resultList, queue)
			}
		}
	} else {
		resultList = queueList
	}

	return &monitor.QueueListResponse{
		QueueNames: resultList,
	}, nil
}

// GetQueueAvailableCores 获取队列及对应的可用核数
func (s *GRPCService) GetQueueAvailableCores(ctx context.Context, req *monitor.GetQueueAvailableCoresRequest) (*monitor.GetQueueAvailableCoresResponse, error) {
	logger := logging.GetLogger(ctx)
	queueInfos, err := s.nodeDao.GetQueueAvailableCores(ctx, req.QueueNames)
	if err != nil {
		logger.Errorf("get queue available cores error, err: %v", err)
		return nil, err
	}

	//过滤队列，排除配置文件中的隐藏队列
	hiddenQueueMap := getHiddenQueue()
	var resultQueues []*monitor.QueueCore
	if len(hiddenQueueMap) != 0 {
		for _, queueInfo := range queueInfos {
			if _, ok := hiddenQueueMap[queueInfo.QueueName]; !ok {
				resultQueues = append(resultQueues, queueInfo)
			}
		}
	} else {
		resultQueues = queueInfos
	}

	return &monitor.GetQueueAvailableCoresResponse{
		QueueCores: resultQueues,
	}, nil
}

// GetQueueCoreInfos 获取 Queue 及核数
func (s *GRPCService) GetQueueCoreInfos(ctx context.Context, req *monitor.GetQueueCoreInfosRequest) (*monitor.GetQueueCoreInfosResponse, error) {
	logger := logging.GetLogger(ctx)
	queueInfos, err := s.nodeDao.GetQueueCoreInfos(ctx)
	if err != nil {
		logger.Errorf("get queue core infos err: %v", err)
		return nil, err
	}

	//过滤队列，排除配置文件中的隐藏队列
	hiddenQueueMap := getHiddenQueue()
	resultQueues := make([]*monitor.QueueCoreInfo, 0)
	if len(hiddenQueueMap) != 0 {
		for _, queueInfo := range queueInfos {
			if _, ok := hiddenQueueMap[queueInfo.QueueName]; !ok {
				resultQueues = append(resultQueues, queueInfo)
			}
		}
	} else {
		resultQueues = queueInfos
	}

	return &monitor.GetQueueCoreInfosResponse{
		QueueCoreInfos: resultQueues,
	}, nil
}

func getHiddenQueue() map[string]bool {
	hiddenQueue := config.GetConfig().Scheduler.HiddenQueue
	hiddenQueuesMap := make(map[string]bool)
	if strutil.IsNotEmpty(hiddenQueue) {
		queues := strings.Split(hiddenQueue, ",")
		for _, queue := range queues {
			hiddenQueuesMap[queue] = true
		}
	}
	return hiddenQueuesMap
}
