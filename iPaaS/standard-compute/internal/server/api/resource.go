package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"

	"github.com/yuansuan/ticp/common/project-root-api/hpc/jobstate"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao/models"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/response"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/util"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/errorcode"
)

const (
	queueKey = "queue"
)

func GetResource(c *gin.Context) {
	logger := trace.GetLogger(c)

	s, err := getState(c)
	if err = response.InternalErrorIfError(c, err, errorcode.InternalServerError); err != nil {
		err = fmt.Errorf("get state from gin ctx failed, %w", err)
		logger.Error(err)
		return
	}

	queues := append(s.Conf.BackendProvider.SchedulerCommon.CandidateQueues, s.Conf.BackendProvider.SchedulerCommon.DefaultQueue)
	if c.Query(queueKey) != "" {
		queues = []string{c.Query(queueKey)}
	}

	resourceInfo, err := s.JobScheduler.GetFreeResource(c, queues)
	if err = response.InternalErrorIfError(c, err, errorcode.InternalServerError); err != nil {
		err = fmt.Errorf("get free resource failed, %w", err)
		logger.Error(err)
		return
	}

	// 实际空闲资源 = 调度器查询的资源 - 标准计算准备中作业占用的资源
	ctx := c.Request.Context()
	jobs, err := dao.Default.GetAllJobsByState(ctx, jobstate.Preparing)
	if err = response.InternalErrorIfError(c, err, errorcode.InternalServerError); err != nil {
		err = fmt.Errorf("get all jobs state = %s failed, %w", jobstate.Preparing, err)
		logger.Error(err)
		return
	}

	calculateResource(jobs, resourceInfo, s.Conf.BackendProvider.SchedulerCommon.CoresPerNode, s.Conf.BackendProvider.SchedulerCommon.ReservedCores)

	response.OK(c, resourceInfo)
}

func calculateResource(jobs []*models.Job, res map[string]*v20230530.Resource, coresPerNode map[string]int, reservedCores map[string]int) {
	for queue := range res {
		res[queue].CoresPerNode = int64(coresPerNode[queue])
		res[queue].ReservedCores = int64(reservedCores[queue])
	}

	if len(jobs) == 0 {
		return
	}
	// jobs group by queue, map[queue]jobs
	jobsMap := make(map[string][]*models.Job)
	for _, job := range jobs {
		jobsMap[job.Queue] = append(jobsMap[job.Queue], job)
	}

	// for each res, calculate the occupied resource
	for queue, jobs := range jobsMap {
		var totalOccupiedNodesNum, totalOccupiedMems int64
		for _, job := range jobs {
			totalOccupiedNodesNum += int64(util.OccupiedNodesNum(int(job.RequestCores), coresPerNode[job.Queue]))
			totalOccupiedMems += job.RequestMemory
		}
		res[queue].Cpu = res[queue].Cpu - totalOccupiedNodesNum*int64(coresPerNode[queue])
		res[queue].IdleNodeNum = res[queue].IdleNodeNum - totalOccupiedNodesNum
		res[queue].Memory = res[queue].Memory - totalOccupiedMems
	}
}
