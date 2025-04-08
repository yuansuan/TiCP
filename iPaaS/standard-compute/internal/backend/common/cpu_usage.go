package common

import (
	"context"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao/models"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

func GetCpuUsage(ctx context.Context, j *models.Job, nodes []string, adjustFactor float64) (*v20230530.CpuUsage, error) {
	log.Infof("GetCpuUsage - Workspace: %s\n", j.Workspace)
	log.Infof("j.ExecHosts: %s\n", j.ExecHosts)

	// 创建一个带有 5 秒超时的新 context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel() // 确保在函数退出时调用 cancel

	var wg sync.WaitGroup
	nodeUsages := make(map[string]float64)
	var totalUsage float64
	var mu sync.Mutex

	done := make(chan struct{})

	go func() {
		for _, node := range nodes {
			wg.Add(1)
			go func(node string) {
				defer wg.Done()
				cmd := exec.CommandContext(ctx, "ssh", node, "top -b -n 1 | grep 'Cpu(s)' | awk -F',' '{print $1}' | awk '{print $2}'")
				output, err := cmd.CombinedOutput()
				log.Debugf("stdout: %s", string(output))
				log.Debugf("stderr: %s", err)
				if err != nil {
					log.Warnf("Error getting CPU usage for node %s: %v\n", node, err)
					return
				}
				cpuUsage, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
				if err != nil {
					log.Warnf("Error parsing CPU usage for node %s: %v\n", node, err)
					return
				}
				log.Infof("Job ID: %v, Node %s current CPU usage: %.2f%%", j.Id, node, cpuUsage)
				mu.Lock()
				nodeUsages[node] = cpuUsage
				totalUsage += cpuUsage
				mu.Unlock()
			}(node)
		}
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Info("All goroutines exited successfully")
	case <-ctx.Done():
		log.Warn("GetCpuUsage timed out after 5 seconds")
	}

	if len(nodeUsages) < len(nodes) {
		log.Warn("Not all nodes reported CPU usage")
		return nil, ErrWrongCPUUsage // 查的不全认为查询失败，外部返回 503
	}

	averageUsage := totalUsage / float64(len(nodes))

	var realCpuUsage float64
	if adjustFactor == 1 {
		realCpuUsage = averageUsage
	} else {
		realCpuUsage = math.Min(averageUsage/adjustFactor, 100.0)
		for node, usage := range nodeUsages {
			nodeUsages[node] = math.Min(usage/adjustFactor, 100.0)
		}
	}

	return &v20230530.CpuUsage{
		JobID:           j.JobID(),
		AverageCpuUsage: realCpuUsage,
		NodeUsages:      nodeUsages,
	}, nil
}
