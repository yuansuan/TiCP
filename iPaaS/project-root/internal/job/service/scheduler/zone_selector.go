package scheduler

import (
	"context"
	"fmt"
	"math"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/resource"
	"github.com/yuansuan/ticp/common/project-root-api/proto/license"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
)

// github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/scheduler ResourceGetter
// go:generate mockgen -destination mock_resource_getter.go -package scheduler
type ResourceGetter interface {
	GetResource(url string) (*resource.SystemGetResponse, error)
}

// ZoneSelector 分区选择器
type ZoneSelector struct {
	zones     schema.Zones
	rsgetter  ResourceGetter
	zoneQueue map[string]string // 应用指定队列
	jobZone   string            // 作业指定分区
	jobQueue  string            // 作业指定队列

	Filters          []ResourceFilter
	OptimalSelectors []ResourceOptimalSelector
}

func NewZoneSelector(zones schema.Zones, rsgetter ResourceGetter,
	zoneQueue map[string]string, jobZone, jobQueue string) *ZoneSelector {
	return &ZoneSelector{
		zones:     zones,
		rsgetter:  rsgetter,
		zoneQueue: zoneQueue,
		jobZone:   jobZone,
		jobQueue:  jobQueue,

		Filters:          make([]ResourceFilter, 0),
		OptimalSelectors: make([]ResourceOptimalSelector, 0),
	}
}

// ResourceFilter 过滤器
type ResourceFilter interface {
	Name() string
	filter(ctx context.Context, zr *util.ZoneResource) (bool, error) // 过滤函数，返回是否通过过滤
}

// ResourceOptimalSelector 优选器
type ResourceOptimalSelector interface {
	Name() string
	optimalSelect(ctx context.Context, zr *util.ZoneResource) (float64, error) // 优选函数，返回一个优先级
	SetWeight(weight float64)
}

func (zs *ZoneSelector) RegisterFilter(filters ...ResourceFilter) {
	zs.Filters = append(zs.Filters, filters...)
}

func (zs *ZoneSelector) RegisterOptimalSelector(optimals ...ResourceOptimalSelector) {
	zs.OptimalSelectors = append(zs.OptimalSelectors, optimals...)
}

// getZoneResource 获取单个分区资源
func (zs *ZoneSelector) getZoneResource(ctx context.Context, zone string) ([]*util.ZoneResource, error) {
	logger := logging.GetLogger(ctx).With("func", "job.ZoneSelector.getZoneResource")
	logger.Debugf("getZoneResource start...")
	defer logger.Debugf("getZoneResource end")

	zoneInfo, ok := zs.zones[zone]
	if !ok {
		logger.Warnf("get zone config failed! zone: %s", zone)
		return nil, fmt.Errorf("get zone config failed! zone: %s", zone)
	}
	if zoneInfo.HPCEndpoint == "" { // 仅有存储的分区，应该属于正常情况
		logger.Infof("zone %s hpc endpoint is empty", zone)
		return nil, nil
	}
	hpcResp, err := zs.rsgetter.GetResource(zoneInfo.HPCEndpoint)
	if err != nil {
		logger.Warnf("get hpc resource error! zone: %s, err: %v", zone, err)
		return nil, err
	}
	if hpcResp == nil || hpcResp.Data == nil {
		logger.Warnf("get hpc resource failed! zone: %s", zone)
		return nil, fmt.Errorf("get hpc resource failed! zone: %s", zone)
	}

	convertFunc := func(q string, rs *schema.Resource) (*util.ZoneResource, error) {
		if rs.CoresPerNode == 0 {
			logger.Errorf("hpc resource coresPerNode is 0!")
			return nil, fmt.Errorf("hpc resource coresPerNode is 0!")
		}
		return &util.ZoneResource{
			Zone:          zone,
			Queue:         q,
			CPU:           rs.Cpu,
			Mem:           rs.Memory,
			CoresPerNode:  rs.CoresPerNode,
			ReservedCores: rs.ReservedCores,
			IsDefault:     rs.IsDefault,
			TotalNodeNum:  rs.TotalNodeNum,
		}, nil
	}

	zoneResources := make([]*util.ZoneResource, 0)
	// 作业指定队列存在时优先级最高
	if zs.jobQueue != "" {
		if hpcResource, ok := hpcResp.Data[zs.jobQueue]; ok {
			zoneResource, err := convertFunc(zs.jobQueue, hpcResource)
			if err != nil {
				return nil, err
			}
			zoneResources = append(zoneResources, zoneResource)
			return zoneResources, nil
		}
	}
	// 应用指定队列中存在该分区，取应用指定队列
	if queue, exists := zs.zoneQueue[zone]; exists {
		if hpcResource, ok := hpcResp.Data[queue]; ok {
			zoneResource, err := convertFunc(queue, hpcResource)
			if err != nil {
				return nil, err
			}
			zoneResources = append(zoneResources, zoneResource)
			return zoneResources, nil
		}
	}
	// 其他队列
	for k, resource := range hpcResp.Data {
		zoneResource, err := convertFunc(k, resource)
		if err != nil {
			return nil, err
		}
		zoneResources = append(zoneResources, zoneResource)
	}
	return zoneResources, nil
}

// ResourceAggregation 获取所有分区资源的聚合
func (zs *ZoneSelector) ResourceAggregation(ctx context.Context) ([]*util.ZoneResource, error) {
	logger := logging.GetLogger(ctx).With("func", "job.ZoneSelector.resourceAggregation")
	logger.Debugf("ResourceAggregation() start...")
	defer logger.Debugf("ResourceAggregation() end")

	resources := make([]*util.ZoneResource, 0)
	for zone := range zs.zones {
		if zs.jobZone != "" && zone != zs.jobZone {
			// 已选取分区情况下，只获取已选取分区的资源即可
			continue
		}
		// 获取分区资源，每个分区可能会返回多个队列的资源
		zrs, err := zs.getZoneResource(ctx, zone)
		if err != nil {
			if zs.isSpecifiedZone(zone) { // 如果指定的分区资源聚合失败，直接返回错误
				return resources, fmt.Errorf("get job zone %s resource error! err: %v", zone, err)
			}
			continue
		}
		if zrs != nil {
			resources = append(resources, zrs...)
		}
	}
	return resources, nil
}

type SelectParams struct {
	Resources     []*util.ZoneResource
	JobID         snowflake.ID
	PreScheduleID snowflake.ID
	UserID        snowflake.ID
}

type SelectOption func(context.Context, *SelectParams) context.Context

func WithJobID(jobID snowflake.ID) SelectOption {
	return func(ctx context.Context, params *SelectParams) context.Context {
		ctx = logging.AppendWith(ctx, "job_id", jobID)
		params.JobID = jobID
		return ctx
	}
}

func WithPreScheduleID(preScheduleID snowflake.ID) SelectOption {
	return func(ctx context.Context, params *SelectParams) context.Context {
		ctx = logging.AppendWith(ctx, "pre_schedule_id", preScheduleID)
		params.PreScheduleID = preScheduleID
		return ctx
	}
}

func WithUserID(userID snowflake.ID) SelectOption {
	return func(ctx context.Context, params *SelectParams) context.Context {
		ctx = logging.AppendWith(ctx, "user_id", userID)
		params.UserID = userID
		return ctx
	}
}

func (zs *ZoneSelector) Select(ctx context.Context, params *SelectParams,
	opts ...SelectOption) (string, *util.ZoneResource, error) {
	ctx = logging.AppendWith(ctx, "func", "job.ZoneSelector.SelectNew")
	for _, opt := range opts {
		ctx = opt(ctx, params)
	}
	logger := logging.GetLogger(ctx)
	logger.Debugf("app queueMap is %v", zs.zoneQueue)
	logger.Debugf("SelectNew start...")
	defer logger.Debugf("SelectNew end")
	// 过滤, 过滤掉不可用的资源
	allowedZones, err := zs.Filter(ctx, params.Resources)
	if err != nil {
		return "", nil, err
	}
	// 优选, 可用资源中选择一个最合适的资源
	optimalZone, err := zs.OptimalSelect(ctx, allowedZones)
	if err != nil {
		return "", nil, err
	}
	return optimalZone.Zone, optimalZone, nil
}

// Filter 过滤掉不可用的资源,返回为allowedZones,blockedZones,error
func (zs *ZoneSelector) Filter(ctx context.Context, zones []*util.ZoneResource) ([]*util.ZoneResource, error) {
	logger := logging.GetLogger(ctx).With("func", "job.ZoneSelector.filter")
	logger.Debugf("filter start...")
	defer logger.Debugf("filter end")

	if len(zones) == 0 {
		logger.Warnf("no zone resources")
		return nil, fmt.Errorf("no zone resources")
	}

	allowedZones := make([]*util.ZoneResource, 0)
	for _, zr := range zones {
		if passed, err := zs.passesAllFilters(ctx, zr); err != nil {
			logger.Infof("filter zone %s error! err: %v", zr.Zone, err)
			if zs.isSpecifiedZone(zr.Zone) && (zs.isSpecifiedQueue(zr) || len(zones) == 1) {
				// 指定队列或分区只有一个队列时，直接返回错误
				return allowedZones, fmt.Errorf("job zone %s queue %s was filtered out,meet err: %v",
					zr.Zone, zr.Queue, err)
			}
		} else if !passed {
			logger.Infof(zr.PrintFilterReason())
			if zs.isSpecifiedZone(zr.Zone) && (zs.isSpecifiedQueue(zr) || len(zones) == 1) {
				// 指定队列或分区只有一个队列时，直接返回错误
				return allowedZones, fmt.Errorf("job zone %s queue %s was filtered out, reason: %s",
					zr.Zone, zr.Queue, zr.PrintFilterReason())
			}
		} else {
			allowedZones = append(allowedZones, zr)
		}
	}
	allowedCount := len(allowedZones)
	logger.Infof("%d zones were filtered out, %d zones remain.", len(zones)-allowedCount, allowedCount)
	if len(allowedZones) == 0 {
		logger.Infof("all zones were filtered out")
		return allowedZones, fmt.Errorf("all zones were filtered out")
	}
	return allowedZones, nil
}

func (zs *ZoneSelector) isSpecifiedZone(zone string) bool {
	return zs.jobZone != "" && zs.jobZone == zone
}

func (zs *ZoneSelector) isSpecifiedQueue(zr *util.ZoneResource) bool {
	if zs.jobQueue != "" && zs.jobQueue == zr.Queue {
		return true
	}
	if queue, exists := zs.zoneQueue[zr.Zone]; exists && queue == zr.Queue {
		return true
	}
	return false
}

func (zs *ZoneSelector) passesAllFilters(ctx context.Context, zr *util.ZoneResource) (bool, error) {
	logger := logging.GetLogger(ctx).With("func", "job.ZoneSelector.passesAllFilters", "current_zone", zr.Zone)
	for _, f := range zs.Filters {
		if filterPassed, err := f.filter(ctx, zr); err != nil {
			return false, err
		} else if !filterPassed {
			logger.Debugf("filter %s not passed", f.Name())
			return false, nil
		}
	}
	return true, nil
}

// OptimalSelect 优选一个最合适的资源
func (zs *ZoneSelector) OptimalSelect(ctx context.Context, zones []*util.ZoneResource) (*util.ZoneResource, error) {
	logger := logging.GetLogger(ctx).With("func", "job.ZoneSelector.optimalSelect")
	logger.Debugf("optimalSelect start...")
	defer logger.Debugf("optimalSelect end")

	var selectedZone *util.ZoneResource
	var maxPriority float64 = math.Inf(-1)
	for _, zr := range zones {
		priority, err := zs.calculatePriority(ctx, zr, zs.OptimalSelectors) // 计算优先级
		if err != nil {
			logger.Warnf("zone %s calculate priority error! err: %v", zr.Zone, err)
			continue
		}
		logger.Infof("zone %s priority: %f", zr.Zone, priority)
		if priority > maxPriority {
			maxPriority = priority
			selectedZone = zr
		}
	}

	if selectedZone == nil {
		logger.Warnf("no optimal zone selected")
		return nil, fmt.Errorf("no optimal zone selected")
	}

	logger.Infof("optimal zone selected, zone: %s, resource Cpu: %d, Mem: %d, PerNode: %d",
		selectedZone.Zone, selectedZone.CPU, selectedZone.Mem, selectedZone.CoresPerNode)
	return selectedZone, nil
}

// calculatePriority 计算优先级
func (zs *ZoneSelector) calculatePriority(ctx context.Context, zone *util.ZoneResource,
	rs []ResourceOptimalSelector) (float64, error) {
	priority := float64(0)
	for _, r := range rs {
		p, err := r.optimalSelect(ctx, zone)
		if err != nil {
			return math.Inf(-1), err
		}
		priority += p
	}
	return priority, nil
}

// baseFilter 基础资源过滤器
type baseFilter struct {
	CPURange     util.CoresRange
	Mem          int64
	Zone         string
	Queue        string
	InputHPCZone string
	Shared       bool
}

// NewBaseFilter 检查资源是否满足基本要求的过滤器
func NewBaseFilter(cpus util.CoresRange, mem int64, zone string, queue string,
	inputHPCZone string, shared bool) *baseFilter {
	return &baseFilter{
		CPURange:     cpus,
		Mem:          mem,
		Zone:         zone,
		Queue:        queue,
		InputHPCZone: inputHPCZone,
		Shared:       shared,
	}
}

func (b *baseFilter) Name() string {
	return "BaseFilter"
}

func (b *baseFilter) filter(ctx context.Context, zr *util.ZoneResource) (bool, error) {
	for _, fn := range []func(context.Context, *util.ZoneResource) (bool, error){
		b.filterCPU, b.filterMem, b.filterZone, b.filterQueue, b.filterInputHPCZone} {
		passed, err := fn(ctx, zr)
		if err != nil {
			return false, err
		}
		if !passed {
			return false, nil
		}
	}
	return true, nil
}

func (b *baseFilter) filterCPU(ctx context.Context, zr *util.ZoneResource) (bool, error) {
	_, resourceUsageCpus, err := util.CalculateResourceUsage(b.CPURange,
		zr.CoresPerNode, zr.CPU, b.Shared)
	if err != nil {
		zr.AddFilterReason(fmt.Sprintf("calculate resource usage error: %v", err))
		return false, err
	}

	// 过滤cpu不够的
	if zr.CPU < resourceUsageCpus {
		zr.AddFilterReason(fmt.Sprintf("cpu not enough, need cpu: %d, zone cpu: %d",
			resourceUsageCpus, zr.CPU))
		return false, nil
	}

	return true, nil
}

func (b *baseFilter) filterMem(ctx context.Context, zr *util.ZoneResource) (bool, error) {
	// 过滤mem不够的
	if zr.Mem < b.Mem {
		zr.AddFilterReason(fmt.Sprintf("mem not enough, need mem: %d, zone mem: %d", b.Mem, zr.Mem))
		return false, nil
	}
	return true, nil
}

func (b *baseFilter) filterZone(ctx context.Context, zr *util.ZoneResource) (bool, error) {
	// 如果指定了分区，过滤掉不是指定分区的
	if len(b.Zone) > 0 && zr.Zone != b.Zone {
		zr.AddFilterReason(fmt.Sprintf("zone not match, need zone: %s, "+
			"currently selected zone: %s", b.Zone, zr.Zone))
		return false, nil
	}
	return true, nil
}

func (b *baseFilter) filterQueue(ctx context.Context, zr *util.ZoneResource) (bool, error) {
	// 作业有指定队列时，过滤掉不是指定队列的
	if len(b.Queue) > 0 && zr.Queue != b.Queue {
		zr.AddFilterReason(fmt.Sprintf("queue not match, need queue: %s, "+
			"currently selected queue: %s", b.Queue, zr.Queue))
		return false, nil
	}
	return true, nil
}

func (b *baseFilter) filterInputHPCZone(ctx context.Context, zr *util.ZoneResource) (bool, error) {
	// input为HPC时，过滤掉不是inputHPCZone的
	if len(b.InputHPCZone) > 0 && zr.Zone != b.InputHPCZone {
		zr.AddFilterReason(fmt.Sprintf("inputHPCZone not match, InputHPCZone: %s, "+
			"currently selected zone: %s", b.InputHPCZone, zr.Zone))
		return false, nil
	}
	return true, nil
}

type licenseFilter struct {
	zones  schema.Zones
	params *LicenseFilterParams
	app    *models.Application

	licenseClient license.LicenseManagerServiceClient
}

type LicenseFilterParams struct {
	IdentifierID         snowflake.ID // jobID或prescheduleID
	CPURange             util.CoresRange
	JobResourceUsageCpus int64 // 仅average模式下用这个字段
	Shared               bool
	Average              bool // 预调度没有这个字段
	Type                 string
}

// NewLicenseFilter 检查分区是否有license的过滤器
func NewLicenseFilter(zones schema.Zones, params *LicenseFilterParams,
	app *models.Application, licenseClient license.LicenseManagerServiceClient) *licenseFilter {
	return &licenseFilter{zones: zones, params: params, app: app, licenseClient: licenseClient}
}

func (lf *licenseFilter) Name() string {
	return "LicenseFilter"
}

func (lf *licenseFilter) filter(ctx context.Context, zr *util.ZoneResource) (bool, error) {
	logger := logging.GetLogger(ctx).With("func", "job.ZoneSelector.licenseFilter.filter")
	// license过滤
	if !config.GetConfig().ChangeLicense || !lf.app.LicManagerId.NotZero() {
		logger.Info("ChangeLicense or licManager not set, means no need to check license")
		return true, nil
	}

	// 先release这个jobID的license
	releaseReq := &license.ReleaseRequest{
		JobId: lf.params.IdentifierID.Int64(),
	}
	_, err := lf.licenseClient.ReleaseLicense(ctx, releaseReq)
	if err != nil {
		logger.Warnf("job %s release License networks: %v", lf.params.IdentifierID.String(), err)
		zr.AddFilterReason(fmt.Sprintf("job %s release License networks: %v", lf.params.IdentifierID.String(), err))
		return false, err
	}
	var totalCores int64
	if lf.params.Average {
		if lf.params.Shared || lf.app.NeedLimitCore {
			totalCores = lf.params.JobResourceUsageCpus
		} else { // 非shared用户至少要用一个单机节点核数
			totalCores = max(zr.CoresPerNode, lf.params.JobResourceUsageCpus)
		}
	} else {
		// 再查询这个jobID的license
		// shared或者app.NeedLimitCore为true时，不取整
		_, resourceUsageCpus, err := util.CalculateResourceUsage(lf.params.CPURange,
			zr.CoresPerNode, zr.CPU, lf.params.Shared || lf.app.NeedLimitCore)
		if err != nil {
			logger.Warnf("calculate resource usage error: %v", err)
			zr.AddFilterReason(fmt.Sprintf("calculate resource usage error: %v", err))
			return false, err
		}
		totalCores = resourceUsageCpus
	}

	consumeInfo := &license.ConsumeInfo{
		JobId:        lf.params.IdentifierID.Int64(),
		AppId:        lf.app.ID.Int64(),
		Cpus:         totalCores,
		LicManagerId: lf.app.LicManagerId.Int64(),
		HpcEndpoint:  lf.zones[zr.Zone].HPCEndpoint,
	}

	// 请求license的RPC
	var consumeInfos []*license.ConsumeInfo
	consumeInfos = append(consumeInfos, consumeInfo)

	acquireReq := &license.ConsumeRequest{
		Info:      consumeInfos,
		OnlyQuery: true, // 仅查询
	}

	licenseServer, err := lf.licenseClient.AcquireLicenses(ctx, acquireReq) //调 license的rpc接口
	if err != nil {
		logger.Warnf("request license server,acquire License networks: error: %v", err)
		zr.AddFilterReason(fmt.Sprintf("request license server,acquire License networks: error: %v", err))
		return false, err
	}

	results := licenseServer.Result
	result := results[0]

	licenseStatus := result.Status
	// license 未配置
	if licenseStatus == license.LicenseStatus_UNCONFIGURED {
		logger.Warnf("license not configured, LicManagerId: [%v]", lf.app.LicManagerId.String())
		return true, nil //! 部分应用未配置
	}

	// license 不够
	if licenseStatus == license.LicenseStatus_NOTENOUTH ||
		licenseStatus == license.LicenseStatus_UNPUBLISH {
		logger.Warnf("lack of license, need cpus %d, LicManagerId %s, state: [%s], "+
			"if long-term alarm, check whether the total number of licenses is insufficient or "+
			"the partition has no license configuration", totalCores, lf.app.LicManagerId.String(),
			licenseStatus.String())
		zr.AddFilterReason(fmt.Sprintf("lack of license, need cpus %d, LicManagerId %s, state: [%s]",
			totalCores, lf.app.LicManagerId.String(), licenseStatus.String()))
		return false, nil
	}
	return true, nil
}

type ReserveResourceFilter struct {
	CPURange  util.CoresRange
	Queue     string
	ZoneQueue map[string]string // 应用指定队列
}

func NewReserveResourceFilter(cpus util.CoresRange, queue string,
	zoneQueue map[string]string) *ReserveResourceFilter {
	return &ReserveResourceFilter{CPURange: cpus, Queue: queue, ZoneQueue: zoneQueue}
}

func (rrf *ReserveResourceFilter) Name() string {
	return "ReserveResourceFilter"
}

func (rrf *ReserveResourceFilter) filter(ctx context.Context, zr *util.ZoneResource) (bool, error) {
	logger := logging.GetLogger(ctx).With("func", "job.ZoneSelector.ReserveResourceFilter.filter")
	logger.Debugf("ReserveResourceFilter start...")
	defer logger.Debugf("ReserveResourceFilter end")

	if rrf.Queue != "" {
		// 指定了队列的，不过滤预留资源，即直接使用
		return true, nil
	}
	// 应用指定队列中存在该分区，不过滤预留资源，即直接使用
	if queue, exists := rrf.ZoneQueue[zr.Zone]; exists {
		if queue == zr.Queue {
			return true, nil
		}
	}
	// 预留资源过滤
	_, resourceUsageCpus, err := util.CalculateResourceUsage(rrf.CPURange, zr.CoresPerNode, zr.CPU, false)
	if err != nil {
		logger.Warnf("calculate resource usage error: %v", err)
		zr.AddFilterReason(fmt.Sprintf("calculate resource usage error: %v", err))
		return false, err
	}
	// 剩余核数减去需求核数不能低于队列预留核数
	if zr.ReservedCores != 0 && zr.CPU-resourceUsageCpus < zr.ReservedCores {
		zr.AddFilterReason(fmt.Sprintf("cpu not enough, need cpus: %d, zone cpus: %d, "+
			"zone reserved cpus: %d", resourceUsageCpus, zr.CPU, zr.ReservedCores))
		return false, nil
	}
	return true, nil
}

// storageFirst 存储优先
type storageFirst struct {
	weight      float64
	storageZone string
}

// NewStorageFirst 存储优先优选器
func NewStorageFirst(storageZone string) *storageFirst {
	return &storageFirst{
		storageZone: storageZone,
	}
}

func (sf *storageFirst) Name() string {
	return "StorageFirst"
}

func (sf *storageFirst) optimalSelect(ctx context.Context, zr *util.ZoneResource) (float64, error) {
	priority := float64(0) // 优先级

	if zr.Zone == sf.storageZone {
		priority = 1
	}

	return priority * sf.weight, nil
}

func (sf *storageFirst) SetWeight(weight float64) {
	sf.weight = weight
}

// 资源数量
type resourceCount struct {
	weight float64
}

// NewResourceCount 资源数量优选器
func NewResourceCount() *resourceCount {
	return &resourceCount{}
}

func (rc *resourceCount) Name() string {
	return "ResourceCount"
}

func (rc *resourceCount) optimalSelect(ctx context.Context, zr *util.ZoneResource) (float64, error) {
	// 资源越多优先级越高
	priority := float64(zr.CPU) / 10000 // 优先级
	return priority * rc.weight, nil
}

func (rc *resourceCount) SetWeight(weight float64) {
	rc.weight = weight
}

// 分区顺序
type zonePrioritySelector struct {
	weight        float64
	zonesPriority map[string]int
}

// NewZonePrioritySelector 分区顺序优选器
func NewZonePrioritySelector(zones []string) *zonePrioritySelector {
	// 根据zones的顺序，给每个分区一个优先级
	zonesPriority := make(map[string]int)
	for i, zone := range zones {
		// 排序从大到小
		zonesPriority[zone] = len(zones) - i
	}

	return &zonePrioritySelector{zonesPriority: zonesPriority}
}

func (zps *zonePrioritySelector) Name() string {
	return "ZonePrioritySelector"
}

func (zps *zonePrioritySelector) optimalSelect(ctx context.Context, zr *util.ZoneResource) (float64, error) {
	priority := float64(zps.zonesPriority[zr.Zone]) // 优先级
	return priority * zps.weight, nil
}

func (zps *zonePrioritySelector) SetWeight(weight float64) {
	zps.weight = weight
}

type queuePrioritySelector struct {
	weight float64
}

// NewQueuePrioritySelector 队列优先级优选器
func NewQueuePrioritySelector() *queuePrioritySelector {
	return &queuePrioritySelector{}
}

func (dqs *queuePrioritySelector) Name() string {
	return "QueuePrioritySelector"
}

func (dqs *queuePrioritySelector) optimalSelect(ctx context.Context, zr *util.ZoneResource) (float64, error) {
	priority := float64(0) // 优先级
	// 可拓展：一个分区不仅有默认队列和另一个队列，两个以上队列也可按优先级选取
	if zr.IsDefault {
		priority = 1
	}
	return priority * dqs.weight, nil
}

func (dqs *queuePrioritySelector) SetWeight(weight float64) {
	dqs.weight = weight
}
