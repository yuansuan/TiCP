package api

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao/models"

	tassert "github.com/stretchr/testify/assert"
)

func Test_calculateResource(t *testing.T) {
	configString := `
backend-provider:
  scheduler-common:
    default-queue: g1_test_1
    candidate-queues: [g1_test_2]
    cores-per-node:
      g1_test_1: 56
      g1_test_2: 24
    reserved-cores:
      g1_test_1: 56
`
	cfg := new(config.Config)
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(strings.NewReader(configString))
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(cfg, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.Metadata = &mapstructure.Metadata{}
	})
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	tc := []struct {
		name string
		jobs []*models.Job
	}{
		{
			name: "empty jobs",
			jobs: []*models.Job{},
		},
		{
			name: "jobs in g1_test_1",
			jobs: []*models.Job{
				{
					Queue:         "g1_test_1",
					RequestCores:  56,
					RequestMemory: 200,
				},
			},
		},
		{
			name: "jobs in g1_test_2",
			jobs: []*models.Job{
				{
					Queue:         "g1_test_2",
					RequestCores:  24,
					RequestMemory: 200,
				},
			},
		},
		{
			name: "jobs in both",
			jobs: []*models.Job{
				{
					Queue:         "g1_test_1",
					RequestCores:  56,
					RequestMemory: 200,
				},
				{
					Queue:         "g1_test_2",
					RequestCores:  24,
					RequestMemory: 200,
				},
			},
		},
	}

	for _, tt := range tc {
		res := map[string]*v20230530.Resource{
			"g1_test_1": {
				Cpu:          4592,
				TotalCpu:     5040,
				IdleNodeNum:  82,
				TotalNodeNum: 90,
				Memory:       16276708,
				TotalMemory:  17141760,
				IsDefault:    true,
			},
			"g1_test_2": {
				Cpu:          560,
				TotalCpu:     560,
				IdleNodeNum:  10,
				TotalNodeNum: 10,
				Memory:       1846123,
				TotalMemory:  1904640,
				IsDefault:    false,
			},
		}

		t.Run(tt.name, func(t *testing.T) {
			calculateResource(tt.jobs, res, cfg.BackendProvider.SchedulerCommon.CoresPerNode, cfg.BackendProvider.SchedulerCommon.ReservedCores)
			tassert.Equal(t, int64(56), res["g1_test_1"].CoresPerNode)
			tassert.Equal(t, int64(24), res["g1_test_2"].CoresPerNode)
			tassert.Equal(t, int64(56), res["g1_test_1"].ReservedCores)
			tassert.Equal(t, int64(0), res["g1_test_2"].ReservedCores)
		})
	}
}
