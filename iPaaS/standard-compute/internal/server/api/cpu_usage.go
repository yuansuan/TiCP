// cpu_usage.go

package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/jobstate"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/cpuusage"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend/common"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/response"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/errorcode"
	"strings"
)

func GetCpuUsage(c *gin.Context) {
	logger := trace.GetLogger(c)

	req := new(cpuusage.SystemGetRequest)
	err := c.ShouldBindUri(req)
	if err = response.BadRequestIfError(c, err, errorcode.InvalidArgument); err != nil {
		err = fmt.Errorf("bind get job cpu usage request failed, %w", err)
		logger.Info(err)
		return
	}

	// 将 job_id 转换为 int64
	jobID, err := snowflake.ParseString(req.JobID)
	if err = response.BadRequestIfError(c, err, errorcode.InvalidJobID); err != nil {
		err = fmt.Errorf("parse job id %s to snowflake id failed, %w", req.JobID, err)
		logger.Info(err)
		return
	}

	// 从数据库获取作业信息
	ctx := c.Request.Context()
	exist, j, err := dao.Default.GetJob(ctx, jobID.Int64())
	if err = response.InternalErrorIfError(c, err, errorcode.DatabaseInternalServerError); err != nil {
		err = fmt.Errorf("[cpu usage]get job from db failed, %w", err)
		logger.Error(err)
		return
	}
	if !exist {
		err = fmt.Errorf("job not found")
		_ = response.NotfoundIfError(c, err, errorcode.JobNotFound)
		logger.Info(err)
		return
	}

	// 只有 running 状态才检查cpu
	if j.State != jobstate.Running {
		err := fmt.Errorf("job ID: %v, %v, not in running state, current state: %v", req.JobID, j.Id, j.State)
		_ = response.ForbiddenIfError(c, err, errorcode.JobNotRunning)
		logger.Warn(err)
		return
	}

	// 先校验job参数再调 get cpu usage
	//1. 没有查到节点，外部返回 500
	if j.ExecHosts == "" {
		err := common.ErrNoNodesAvailable
		_ = response.InternalErrorIfError(c, err, errorcode.InternalServerError)
		logger.Error(fmt.Errorf("get CPU usage failed: %w", err))
		return
	}

	nodes := strings.Split(j.ExecHosts, ",")
	nodesNum := len(nodes)

	// 根据partition name读取单个节点的核数
	severCpusPerNode := int64(config.GetConfig().BackendProvider.SchedulerCommon.CoresPerNode[j.Queue])

	// 2. 根据j.Queue读到的配置文件里 节点的CPU数量为0 或者 读取失败
	// 外部返回500
	if severCpusPerNode == 0 {
		err := common.ErrZeroCPUsPerNode
		_ = response.InternalErrorIfError(c, err, errorcode.InternalServerError)
		logger.Error(fmt.Errorf("get CPU usage failed: %w", err))
		return
	}
	serverCpus := int64(nodesNum) * severCpusPerNode //服务器总cpu核数
	allocCpus := j.AllocCores                        // 给作业分配的总cpu核数
	adjustFactor := float64(allocCpus) / float64(serverCpus)

	// 3.校验adjustFactor
	// 为0说明前面的计算 allocCpus为0；
	// 提前返回 防止后面的计算分母为0
	// 外部返回 500
	if adjustFactor == 0 {
		err := common.ErrZeroAllocCPU
		_ = response.InternalErrorIfError(c, err, errorcode.InternalServerError)
		logger.Error(fmt.Errorf("get CPU usage failed: %w", err))
	}

	s, err := getState(c)
	if err = response.InternalErrorIfError(c, err, errorcode.InternalServerError); err != nil {
		err = fmt.Errorf("get state from gin ctx failed, %w", err)
		logger.Info(err)
		return
	}
	cpuUsage, err := s.JobScheduler.GetCpuUsage(ctx, j, nodes, adjustFactor)

	if err != nil {
		if errors.Is(err, common.ErrWrongCPUUsage) {
			_ = response.ServiceUnavailableIfError(c, err, errorcode.WrongCPUUsage) //返回503
		} else {
			_ = response.InternalErrorIfError(c, err, errorcode.InternalServerError)
		}
		logger.Error(fmt.Errorf("get CPU usage failed: %w", err))
		return
	}

	Data := &v20230530.CpuUsage{
		JobID:           jobID.String(),
		AverageCpuUsage: cpuUsage.AverageCpuUsage,
		NodeUsages:      cpuUsage.NodeUsages,
	}

	response.OK(c, Data)
}
