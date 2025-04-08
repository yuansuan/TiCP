package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigParsing(t *testing.T) {
	// Load the configuration
	cfg, err := NewConfig(WithPath("../config/config.yaml"))
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Test general config values
	assert.Equal(t, "0.0.0.0:8080", cfg.HttpAddress)
	assert.Equal(t, "127.0.0.1:8081", cfg.PerformanceAddress)
	assert.Equal(t, "http://127.0.0.1:8899", cfg.HpcStorageAddress)

	// Test BackendProvider configuration
	assert.NotNil(t, cfg.BackendProvider)
	assert.Equal(t, "slurm", cfg.BackendProvider.Type)
	assert.Equal(t, 3, cfg.BackendProvider.CheckAliveInterval)

	// Test SchedulerCommon configuration
	sc := cfg.BackendProvider.SchedulerCommon
	assert.Equal(t, "Job1", sc.DefaultQueue)
	assert.Equal(t, []string{"job2"}, sc.CandidateQueues)
	assert.Equal(t, "./jobs", sc.Workspace)

	// Test CoresPerNode map
	expectedCoresPerNode := map[string]int{
		"Job1": 20,
		"job2": 25,
	}
	assert.Equal(t, expectedCoresPerNode, sc.CoresPerNode)

	// Test ReservedCores map
	expectedReservedCores := map[string]int{
		"Job1": 20,
	}
	assert.Equal(t, expectedReservedCores, sc.ReservedCores)

	t.Run("Test CoresPerNode usage with Job struct", func(t *testing.T) {
		jobs := []MockJob{
			{Queue: "Job1"},
			{Queue: "job2"},
			{Queue: "non_existent_queue"},
		}

		for i := range jobs {
			j := &jobs[i]
			j.CoresPerNode = int64(cfg.BackendProvider.SchedulerCommon.CoresPerNode[j.Queue])
		}

		assert.Equal(t, int64(20), jobs[0].CoresPerNode, "CoresPerNode for Job1 should be 20")
		assert.Equal(t, int64(25), jobs[1].CoresPerNode, "CoresPerNode for Job2 should be 25")
		assert.Equal(t, int64(0), jobs[2].CoresPerNode, "CoresPerNode for non-existent queue should be 0")
	})

	// Test submit user configuration
	assert.Equal(t, "yuansuan", sc.SubmitSysUser)
	assert.Equal(t, 0, sc.SubmitSysUserUid)
	assert.Equal(t, 0, sc.SubmitSysUserGid)

	// Test Slurm configuration
	assert.NotNil(t, cfg.BackendProvider.Slurm)
	assert.Contains(t, cfg.BackendProvider.Slurm.Submit, "sbatch")
	assert.Equal(t, "scancel ${job_id}", cfg.BackendProvider.Slurm.Kill)

	// Test Singularity configuration
	assert.NotNil(t, cfg.Singularity)
	assert.Equal(t, "/shared/singularity", cfg.Singularity.Storage)
	assert.NotNil(t, cfg.Singularity.Registry)
	assert.Equal(t, "standard-compute-1252829527", cfg.Singularity.Registry.Bucket)

	// Test Snowflake configuration
	assert.Equal(t, int64(1), cfg.Snowflake.Node)

	// Test PreparedFilePath
	assert.Equal(t, "/tmp", cfg.PreparedFilePath)
}

// 模拟作业结构，防止直接引用造成循环依赖
type MockJob struct {
	Queue        string
	CoresPerNode int64
}
