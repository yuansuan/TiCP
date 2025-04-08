package scheduler

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	resource "github.com/yuansuan/ticp/common/project-root-api/hpc/v1/resource"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

type ZoneSelectorSuite struct {
	suite.Suite

	ctrl *gomock.Controller
	w    *httptest.ResponseRecorder
	ctx  *gin.Context

	zones     map[string]*schema.Zone
	zoneQueue map[string]string
	resources []*util.ZoneResource

	*MockResourceGetter
}

func (suite *ZoneSelectorSuite) SetupTest() {

	suite.ctrl = gomock.NewController(suite.T())
	suite.w = httptest.NewRecorder()
	suite.ctx = mock.GinContext(suite.w)

	logger, err := logging.NewLogger(logging.WithDefaultLogConfigOption(), logging.WithUseConsole(true), logging.WithLogLevel("debug"))
	if !suite.NoError(err) {
		return
	}

	suite.ctx.Set(logging.LoggerName, logger)

	suite.MockResourceGetter = NewMockResourceGetter(suite.ctrl)

	suite.zoneQueue = map[string]string{
		"az-jinan":    "q1",
		"az-wuxi":     "qwuxi",
		"az-shanghai": "qsh",
	}

	suite.zones = map[string]*schema.Zone{
		"az-jinan": {
			HPCEndpoint: "https://jn_hpc_endpoint:8080",
		},
		"az-wuxi": {
			HPCEndpoint: "https://wx_hpc_endpoint:8080",
		},
		"az-zhigu": {
			HPCEndpoint: "https://zg_hpc_endpoint:8080",
		},
	}

	suite.MockResourceGetter.EXPECT().GetResource("https://jn_hpc_endpoint:8080").AnyTimes().Return(&resource.SystemGetResponse{
		Response: schema.Response{},
		Data: map[string]*schema.Resource{
			"q1": {
				Cpu:          80,
				CoresPerNode: 10,
				Memory:       200,
				IsDefault:    true,
			},
			"q2": {
				Cpu:          120,
				CoresPerNode: 20,
				Memory:       100,
			},
		},
	}, nil)
	suite.MockResourceGetter.EXPECT().GetResource("https://wx_hpc_endpoint:8080").AnyTimes().Return(&resource.SystemGetResponse{
		Response: schema.Response{},
		Data: map[string]*schema.Resource{
			"qwuxi": {
				Cpu:          120,
				CoresPerNode: 20,
				Memory:       100,
			},
			"qwuxi-def": {
				Cpu:          80,
				CoresPerNode: 10,
				Memory:       200,
				IsDefault:    true,
			},
		},
	}, nil)
	suite.MockResourceGetter.EXPECT().GetResource("https://zg_hpc_endpoint:8080").AnyTimes().Return(&resource.SystemGetResponse{
		Response: schema.Response{},
		Data: map[string]*schema.Resource{
			"qzg": {
				Cpu:          40,
				CoresPerNode: 20,
				Memory:       100,
			},
			"qzg2": {
				Cpu:          20,
				CoresPerNode: 20,
				Memory:       100,
				IsDefault:    true,
			},
		},
	}, nil)
	suite.MockResourceGetter.EXPECT().GetResource("https://sh_hpc_endpoint:8080").AnyTimes().Return(&resource.SystemGetResponse{
		Response: schema.Response{},
		Data: map[string]*schema.Resource{
			"q1": {Cpu: 10, CoresPerNode: 20, Memory: 100, IsDefault: true},
		},
	}, nil)
	suite.MockResourceGetter.EXPECT().GetResource("https://missingDefaultQueue:8080").AnyTimes().Return(&resource.SystemGetResponse{
		Response: schema.Response{},
		Data: map[string]*schema.Resource{
			"someq": {Cpu: 10, CoresPerNode: 20, Memory: 100},
		},
	}, nil)
	suite.MockResourceGetter.EXPECT().GetResource("https://getResourceError:8080").AnyTimes().Return(nil, fmt.Errorf("get resource error"))
	suite.MockResourceGetter.EXPECT().GetResource("https://getResourceNil:8080").AnyTimes().Return(nil, nil)
	suite.MockResourceGetter.EXPECT().GetResource("https://zeroCoresPerNode:8080").AnyTimes().Return(&resource.SystemGetResponse{
		Response: schema.Response{},
		Data: map[string]*schema.Resource{
			"someq2": {Cpu: 10, CoresPerNode: 0, Memory: 100, IsDefault: true},
		},
	}, nil)
}

func (suite *ZoneSelectorSuite) TestGetZoneResource() {
	testCases := []struct {
		name        string
		zone        string
		jobQueue    string
		expectQueue []string
		expectError bool
		mockZones   func() map[string]*schema.Zone
	}{
		{
			name:        "jinan",
			zone:        "az-jinan",
			expectQueue: []string{"q1"},
		},
		{
			name:        "wuxi",
			zone:        "az-wuxi",
			expectQueue: []string{"qwuxi"}, // app.SpecifyQueue的指定队列
		},
		{
			name:        "wuxi-queue",
			zone:        "az-wuxi",
			jobQueue:    "qwuxi-def",           // job指定队列
			expectQueue: []string{"qwuxi-def"}, // job指定队列
		},
		{
			name:        "zhigu",
			zone:        "az-zhigu",
			expectQueue: []string{"qzg", "qzg2"}, // app.SpecifyQueue未指定该分区，默认队列和普通队列都会被返回
		},
		{
			name:        "error zone",
			zone:        "az-error",
			expectError: true,
		},
		{
			name: "queue resource not found", // app.SpecifyQueue指定了队列但实际这个队列不存在, 默认队列
			zone: "az-shanghai",
			mockZones: func() map[string]*schema.Zone {
				return map[string]*schema.Zone{
					"az-shanghai": {
						HPCEndpoint: "https://sh_hpc_endpoint:8080",
					},
				}
			},
			expectQueue: []string{"q1"},
		},
		{
			// 分区无默认队列
			name: "missing default queue",
			zone: "az-missingDefaultQueue",
			// expectError: true,
			mockZones: func() map[string]*schema.Zone {
				return map[string]*schema.Zone{
					"az-missingDefaultQueue": {
						HPCEndpoint: "https://missingDefaultQueue:8080",
					},
				}
			},
			expectQueue: []string{"someq"},
		},
		{
			name:        "az-missingHpcEndpoint",
			zone:        "az-missingHpcEndpoint",
			expectError: false,
			mockZones: func() map[string]*schema.Zone {
				return map[string]*schema.Zone{
					"az-missingHpcEndpoint": {
						HPCEndpoint: "",
					},
				}
			},
		},
		{
			name:        "GetResource error",
			zone:        "az-getResourceError",
			expectError: true,
			mockZones: func() map[string]*schema.Zone {
				return map[string]*schema.Zone{
					"az-getResourceError": {
						HPCEndpoint: "https://getResourceError:8080",
					},
				}
			},
		},
		{
			name:        "GetResource nil",
			zone:        "az-getResourceNil",
			expectError: true,
			mockZones: func() map[string]*schema.Zone {
				return map[string]*schema.Zone{
					"az-getResourceNil": {
						HPCEndpoint: "https://getResourceNil:8080",
					},
				}
			},
		},
		{
			name:        "zeroCoresPerNode error",
			zone:        "az-zeroCoresPerNode",
			expectError: true,
			mockZones: func() map[string]*schema.Zone {
				return map[string]*schema.Zone{
					"az-zeroCoresPerNode": {
						HPCEndpoint: "https://zeroCoresPerNode:8080",
					},
				}
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			zones := suite.zones
			if tc.mockZones != nil {
				zones = tc.mockZones()
			}

			zs := NewZoneSelector(zones, suite.MockResourceGetter, suite.zoneQueue, "", tc.jobQueue)
			zrs, err := zs.getZoneResource(suite.ctx, tc.zone)
			if tc.expectError {
				suite.Error(err)
				return
			}

			if suite.NoError(err) {
				suite.T().Log(zrs)
				if len(zrs) > 0 {
					suite.True(checkResources(tc.expectQueue, zrs))
					// suite.Equal(tc.expectQueue, zrs[0].Queue)
				}
			}
		})
		suite.T().Log("/* -------------------------------------------------------------------------- */")
	}
}

func checkResources(expect []string, zrs []*util.ZoneResource) bool {
	// 使用map来检查所有的zrs中的queue是否都在expect中
	queueMap := make(map[string]bool)
	for _, res := range zrs {
		queueMap[res.Queue] = true
	}

	// 检查expect中的每个元素是否都存在于queueMap中
	for _, item := range expect {
		if !queueMap[item] {
			return false
		}
	}

	// 检查zrs中的每个queue是否都在expect中
	for _, res := range zrs {
		found := false
		for _, item := range expect {
			if res.Queue == item {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func (suite *ZoneSelectorSuite) TestSelect() {
	suite.resources = make([]*util.ZoneResource, 0)
	zs := NewZoneSelector(suite.zones, suite.MockResourceGetter, suite.zoneQueue, "", "")
	resources, err := zs.ResourceAggregation(suite.ctx)
	if !suite.NoError(err) {
		return
	}

	suite.T().Logf(spew.Sdump(resources))

	suite.resources = resources

	testCases := []struct {
		name        string
		Job         *models.Job
		expectZone  string
		expectQueue string
		expectError bool
	}{
		{
			name: "with zone and queue",
			Job: &models.Job{
				Zone:                "az-jinan",
				Queue:               "q1",
				HPCJobID:            "542VhUbwggN",
				UserID:              snowflake.Zero(),
				ResourceUsageCpus:   10,
				ResourceUsageMemory: 10,
				InputType:           string(consts.CloudStorage),
			},
			expectZone:  "az-jinan",
			expectQueue: "q1",
		},
		{
			name: "with zone and queue not in resources",
			Job: &models.Job{
				Zone:                "az-jinan",
				Queue:               "q2", // 选择了分区未被聚合的队列资源
				HPCJobID:            "542VhUbwggN",
				UserID:              snowflake.Zero(),
				ResourceUsageCpus:   10,
				ResourceUsageMemory: 10,
				InputType:           string(consts.CloudStorage),
			},
			expectError: true,
		},
		{
			name: "with zone and queue not enough resource",
			Job: &models.Job{
				Zone:                "az-jinan",
				Queue:               "q1",
				HPCJobID:            "542VhUbwggN",
				UserID:              snowflake.Zero(),
				ResourceUsageCpus:   1000,
				ResourceUsageMemory: 1000,
				InputType:           string(consts.CloudStorage),
			},
			expectError: true,
		},
		{
			name: "error zone",
			Job: &models.Job{
				Zone:                "az-error",
				Queue:               "q1",
				HPCJobID:            "542VhUbwggN",
				UserID:              snowflake.Zero(),
				ResourceUsageCpus:   10,
				ResourceUsageMemory: 10,
				InputType:           string(consts.CloudStorage),
			},
			expectZone:  "az-jinan",
			expectQueue: "q1",
			expectError: true,
		},
		{
			name: "error queue",
			Job: &models.Job{
				Zone:                "az-jinan",
				Queue:               "q123",
				HPCJobID:            "542VhUbwggN",
				UserID:              snowflake.Zero(),
				ResourceUsageCpus:   10,
				ResourceUsageMemory: 10,
				InputType:           string(consts.CloudStorage),
			},
			expectZone:  "az-jinan",
			expectQueue: "q1",
			expectError: true,
		},
		{
			name: "normal select",
			Job: &models.Job{
				HPCJobID:            "542VhUbwggN",
				UserID:              snowflake.Zero(),
				ResourceUsageCpus:   10, // 所有分区资源都够
				ResourceUsageMemory: 10,
				InputType:           string(consts.CloudStorage),
			},
			expectZone:  "az-jinan",
			expectQueue: "q1",
		},
		{
			name: "normal select 2",
			Job: &models.Job{
				HPCJobID:            "542VhUbwggN",
				UserID:              snowflake.Zero(),
				ResourceUsageCpus:   100, // 部分分区资源够
				ResourceUsageMemory: 100,
				InputType:           string(consts.CloudStorage),
			},
			expectZone:  "az-wuxi", // 资源最充足的
			expectQueue: "qwuxi",
		},
		{
			name: "not enough resource",
			Job: &models.Job{
				HPCJobID:            "542VhUbwggN",
				UserID:              snowflake.Zero(),
				ResourceUsageCpus:   1000, // 没有分区资源够
				ResourceUsageMemory: 1000,
				InputType:           string(consts.CloudStorage),
			},
			expectError: true,
		},
		{
			name: "storage first",
			Job: &models.Job{
				HPCJobID:             "542VhUbwggN",
				UserID:               snowflake.Zero(),
				ResourceUsageCpus:    10,
				ResourceUsageMemory:  10,
				FileInputStorageZone: "az-zhigu",
				InputType:            string(consts.CloudStorage),
			},
			expectZone:  "az-zhigu",
			expectQueue: "qzg2",
		},
		{
			name: "with zone no queue",
			Job: &models.Job{
				Zone:                "az-zhigu",
				HPCJobID:            "542VhUbwggN",
				UserID:              snowflake.Zero(),
				ResourceUsageCpus:   10,
				ResourceUsageMemory: 10,
				InputType:           string(consts.CloudStorage),
			},
			expectZone:  "az-zhigu",
			expectQueue: "qzg2",
		},
	}

	// 资源聚合
	for _, tc := range testCases {
		// 子测试
		suite.Run(tc.name, func() {
			job := tc.Job
			suite.resources, err = zs.ResourceAggregation(suite.ctx)
			if !suite.NoError(err) {
				return
			}
			zs.RegisterFilter(
				NewBaseFilter(util.CoresRange{
					MinExpectedCores: job.ResourceUsageCpus,
					MaxExpectedCores: job.ResourceUsageCpus,
				}, job.ResourceUsageMemory, job.Zone, job.Queue, func(inputType, inputStorageZone string) string {
					inputHPCZone := ""
					if inputType == string(consts.HpcStorage) {
						inputHPCZone = inputStorageZone
					}
					return inputHPCZone
				}(job.InputType, job.FileInputStorageZone), false),
			)

			rc := NewResourceCount()
			rc.SetWeight(1)
			qps := NewQueuePrioritySelector() // 队列优先级选择器
			qps.SetWeight(10)
			sf := NewStorageFirst(job.FileInputStorageZone)
			sf.SetWeight(100)
			zs.RegisterOptimalSelector(rc, qps, sf)

			params := &SelectParams{
				Resources: suite.resources,
			}
			zone, selectedZone, err := zs.Select(suite.ctx, params)
			suite.T().Log(zone, selectedZone, err)

			if tc.expectError {
				suite.Error(err)
				return
			}

			if suite.NoError(err) {
				if suite.NotEmpty(zone) && suite.NotNil(selectedZone) {
					suite.Equal(tc.expectZone, zone)
					suite.Equal(tc.expectQueue, selectedZone.Queue)
				}
			}
		})
		suite.T().Log("/* -------------------------------------------------------------------------- */")
		zs.Filters = make([]ResourceFilter, 0)
		zs.OptimalSelectors = make([]ResourceOptimalSelector, 0)
	}
}

func TestZoneSelectorSuite(t *testing.T) {
	suite.Run(t, new(ZoneSelectorSuite))
}
