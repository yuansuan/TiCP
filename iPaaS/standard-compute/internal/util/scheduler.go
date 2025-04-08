package util

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend/job"
)

func EnsureNTaskPerNode(cfg config.SchedulerCommon, j *job.Job) int {
	if j.CoresPerNode != 0 {
		return int(j.CoresPerNode)
	}

	return cfg.CoresPerNode[j.Queue]
}

func getEnvFromFile(envFile string) map[string]string {
	env := make(map[string]string)
	data, err := os.ReadFile(envFile)
	if err != nil {
		return env
	}

	for _, line := range strings.Split(string(data), "\n") {
		kv := strings.Split(line, "=")
		if len(kv) == 2 {
			env[kv[0]] = kv[1]
		}
	}

	return env
}

func UpdateByEnv(j *job.Job) {
	logger := logging.GetLogger(context.TODO())
	envs := getEnvFromFile(filepath.Join(j.Workspace, "__env"))

	j.ExecHosts = envs["YS_NODELIST"]
	j.ExecHostNum, _ = strconv.ParseInt(envs["YS_NUM_NODES"], 10, 64)
	if j.ExecHosts == "" {
		logger.Errorf("获取环境变量失败 - job %v, no YS_NODELIST found", j.Id)
		return
	}
	logger.Infof("UpdateByEnv - job %v YS_NODELIST: %s, workspace: %s, env_file: %s", j.Id, envs["YS_NODELIST"], j.Workspace, filepath.Join(j.Workspace, "__env"))
}
