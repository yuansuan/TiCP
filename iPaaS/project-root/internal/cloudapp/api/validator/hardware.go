package validator

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/hardware"
	"github.com/yuansuan/ticp/common/project-root-api/common"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	zonelib "github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/response"
)

const (
	pageOffsetKey                 = "PageOffset"
	defaultPageOffset             = 0
	pageSizeKey                   = "PageSize"
	defaultPageSize               = 1000
	hardwareNameMaxSize           = 64
	hardwareDescMaxSize           = 255
	hardwareInstanceTypeMaxSize   = 64
	hardwareInstanceFamilyMaxSize = 32
	hardwareCpuModelMaxSize       = 255
	hardwareGpuModelMaxSize       = 255
)

func ValidateAPIGetHardwareRequest(req *hardware.APIGetRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("get hardware request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isHardwareIdValid(req.HardwareId, false)
	if err != nil {
		return fmt.Errorf("HardwareId invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func isHardwareIdValid(hardwareId *string, allowEmpty bool) (error, response.ErrorResp) {
	if hardwareId == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *hardwareId == "" {
		err = fmt.Errorf("[HardwareId] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentHardwareId, err.Error())
	}

	_, err = snowflake.ParseString(*hardwareId)
	if err != nil {
		err = fmt.Errorf("parse [HardwareId] \"%s\" to snowflake id failed, %w", *hardwareId, err)
		return err, response.WrapErrorResp(common.InvalidArgumentHardwareId, err.Error())
	}

	return nil, response.ErrorResp{}
}

func ValidateAPIListHardwareRequest(req *hardware.APIListRequest, c *gin.Context) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("list hardware request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	// PageOffset PageSize 特殊处理一下，如果是查询参数 ?PageSize= ，gin的解析会将他解析为PageSize=0，而期望是赋默认值1000
	err, errResp = ensurePageOffset(&req.PageOffset, c)
	if err != nil {
		return fmt.Errorf("PageOffset invalid, %w", err), errResp
	}

	err, errResp = ensurePageSize(&req.PageSize, c)
	if err != nil {
		return fmt.Errorf("PageSize invalid, %w", err), errResp
	}

	err, errResp = isZoneValid(req.Zone, true)
	if err != nil {
		return fmt.Errorf("Zone invalid, %w", err), errResp
	}

	err, errResp = isHardwareNameValid(req.Name, true)
	if err != nil {
		return fmt.Errorf("Hardware Name invalid, %w", err), errResp
	}

	err, errResp = isCpuValid(req.Cpu)
	if err != nil {
		return fmt.Errorf("Cpu invalid, %w", err), errResp
	}

	err, errResp = isMemValid(req.Mem)
	if err != nil {
		return fmt.Errorf("Mem invalid, %w", err), errResp
	}

	err, errResp = isGpuValid(req.Gpu)
	if err != nil {
		return fmt.Errorf("Gpu invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

// pageOffset 无该查询参数时，pageOffset为nil，*pageOffset不允许操作，故用了指针的指针，为空时将其赋默认值
func ensurePageOffset(pageOffset **int, c *gin.Context) (error, response.ErrorResp) {
	pageOffsetStr := c.Query(pageOffsetKey)
	if pageOffsetStr == "" {
		*pageOffset = PInt(defaultPageOffset)
		return nil, response.ErrorResp{}
	}

	pageOffsetValue, err := strconv.Atoi(pageOffsetStr)
	if err != nil {
		err = fmt.Errorf("invalid [PageOffset] \"%s\"", pageOffsetStr)
		return err, response.WrapErrorResp(common.InvalidPageOffset, err.Error())
	}

	if pageOffsetValue < 0 {
		err = fmt.Errorf("[PageOffset] cannot less than 0")
		return err, response.WrapErrorResp(common.InvalidPageOffset, err.Error())
	}

	*pageOffset = PInt(pageOffsetValue)
	return nil, response.ErrorResp{}
}

// pageSize 无该查询参数时，pageSize为nil，*pageSize不允许操作，故用了指针的指针，为空时将其赋默认值
func ensurePageSize(pageSize **int, c *gin.Context) (error, response.ErrorResp) {
	pageSizeStr := c.Query(pageSizeKey)
	if pageSizeStr == "" {
		*pageSize = PInt(defaultPageSize)
		return nil, response.ErrorResp{}
	}

	pageSizeValue, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		err = fmt.Errorf("invalid [PageSize] \"%s\"", pageSizeStr)
		return err, response.WrapErrorResp(common.InvalidPageSize, err.Error())
	}

	if pageSizeValue < 1 || pageSizeValue > 1000 {
		err = fmt.Errorf("[PageSize] should be in 1-1000")
		return err, response.WrapErrorResp(common.InvalidPageSize, err.Error())
	}

	*pageSize = PInt(pageSizeValue)
	return nil, response.ErrorResp{}
}

func isZoneValid(zone *string, allowEmpty bool) (error, response.ErrorResp) {
	if zone == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	zoneParsed := zonelib.Parse(*zone)
	if !allowEmpty && zoneParsed.IsEmpty() {
		err = fmt.Errorf("[Zone] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentZone, err.Error())
	}

	if !zoneParsed.IsValid() {
		err = fmt.Errorf("[Zone] invalid, %s", *zone)
		return err, response.WrapErrorResp(common.InvalidArgumentZone, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isHardwareNameValid(name *string, allowEmpty bool) (error, response.ErrorResp) {
	if name == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *name == "" {
		err = fmt.Errorf("[Name] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentHardwareName, err.Error())
	}

	if len(*name) > hardwareNameMaxSize {
		err = fmt.Errorf("[Name] too long")
		return err, response.WrapErrorResp(common.InvalidArgumentHardwareName, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isCpuValid(cpu *int) (error, response.ErrorResp) {
	if cpu == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if *cpu <= 0 {
		err = fmt.Errorf("[Cpu] cannot less than or equal to zero")
		return err, response.WrapErrorResp(common.InvalidArgumentCpu, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isMemValid(mem *int) (error, response.ErrorResp) {
	if mem == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if *mem <= 0 {
		err = fmt.Errorf("[Mem] cannot less than or equal to zero")
		return err, response.WrapErrorResp(common.InvalidArgumentMem, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isGpuValid(gpu *int) (error, response.ErrorResp) {
	if gpu == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if *gpu < 0 {
		err = fmt.Errorf("[Mem] cannot less than zero")
		return err, response.WrapErrorResp(common.InvalidArgumentGpu, err.Error())
	}

	return nil, response.ErrorResp{}
}

func ValidateAdminPostHardwaresRequest(req *hardware.AdminPostRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("post hardware request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isZoneValid(req.Zone, false)
	if err != nil {
		return fmt.Errorf("Zone invalid, %w", err), errResp
	}

	err, errResp = isHardwareNameValid(req.Name, false)
	if err != nil {
		return fmt.Errorf("Name invalid, %w", err), errResp
	}

	err, errResp = isHardwareDescValid(req.Desc, true)
	if err != nil {
		return fmt.Errorf("Desc invalid, %w", err), errResp
	}

	err, errResp = isInstanceTypeValid(req.InstanceType, false)
	if err != nil {
		return fmt.Errorf("InstanceType invalid, %w", err), errResp
	}

	err, errResp = isInstanceFamilyValid(req.InstanceFamily, true)
	if err != nil {
		return fmt.Errorf("InstanceFamily invalid, %w", err), errResp
	}

	err, errResp = isNetworkValid(req.Network)
	if err != nil {
		return fmt.Errorf("Network invalid, %w", err), errResp
	}

	err, errResp = isCpuValid(req.Cpu)
	if err != nil {
		return fmt.Errorf("Cpu invalid, %w", err), errResp
	}

	err, errResp = isCpuModelValid(req.CpuModel, true)
	if err != nil {
		return fmt.Errorf("CpuModel invalid, %w", err), errResp
	}

	err, errResp = isMemValid(req.Mem)
	if err != nil {
		return fmt.Errorf("Mem invalid, %w", err), errResp
	}

	err, errResp = isGpuValid(req.Gpu)
	if err != nil {
		return fmt.Errorf("Gpu invalid, %w", err), errResp
	}

	err, errResp = isGpuModelValid(req.GpuModel, true)
	if err != nil {
		return fmt.Errorf("GpuModel invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func isHardwareDescValid(desc *string, allowEmpty bool) (error, response.ErrorResp) {
	if desc == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *desc == "" {
		err = fmt.Errorf("[Desc] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidDesc, err.Error())
	}

	if len(*desc) > hardwareDescMaxSize {
		err = fmt.Errorf("[Desc] too long")
		return err, response.WrapErrorResp(common.InvalidDesc, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isInstanceTypeValid(instanceType *string, allowEmpty bool) (error, response.ErrorResp) {
	if instanceType == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *instanceType == "" {
		err = fmt.Errorf("[InstanceType] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentInstanceType, err.Error())
	}

	if len(*instanceType) > hardwareInstanceTypeMaxSize {
		err = fmt.Errorf("[InstanceType] too long")
		return err, response.WrapErrorResp(common.InvalidArgumentInstanceType, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isInstanceFamilyValid(instanceFamily *string, allowEmpty bool) (error, response.ErrorResp) {
	if instanceFamily == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *instanceFamily == "" {
		err = fmt.Errorf("[InstanceFamily] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentInstanceFamily, err.Error())
	}

	if len(*instanceFamily) > hardwareInstanceFamilyMaxSize {
		err = fmt.Errorf("[InstanceFamily] too long")
		return err, response.WrapErrorResp(common.InvalidArgumentInstanceFamily, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isNetworkValid(network *int) (error, response.ErrorResp) {
	if network == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if *network < 0 {
		err = fmt.Errorf("[Network] cannot less than zero")
		return err, response.WrapErrorResp(common.InvalidArgumentNetwork, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isCpuModelValid(cpuModel *string, allowEmpty bool) (error, response.ErrorResp) {
	if cpuModel == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *cpuModel == "" {
		err = fmt.Errorf("[CpuModel] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentCpuModel, err.Error())
	}

	if len(*cpuModel) > hardwareCpuModelMaxSize {
		err = fmt.Errorf("[CpuModel] too long")
		return err, response.WrapErrorResp(common.InvalidArgumentCpuModel, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isGpuModelValid(gpuModel *string, allowEmpty bool) (error, response.ErrorResp) {
	if gpuModel == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *gpuModel == "" {
		err = fmt.Errorf("[GpuModel] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentGpuModel, err.Error())
	}

	if len(*gpuModel) > hardwareGpuModelMaxSize {
		err = fmt.Errorf("[GpuModel] too long")
		return err, response.WrapErrorResp(common.InvalidArgumentGpuModel, err.Error())
	}

	return nil, response.ErrorResp{}
}

func ValidateAdminPatchHardwareRequest(req *hardware.AdminPatchRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("patch hardware request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isHardwareIdValid(req.HardwareId, false)
	if err != nil {
		return fmt.Errorf("HardwareId invalid, %w", err), errResp
	}

	err, errResp = isZoneValid(req.Zone, true)
	if err != nil {
		return fmt.Errorf("Zone invalid, %w", err), errResp
	}

	err, errResp = isHardwareNameValid(req.Name, true)
	if err != nil {
		return fmt.Errorf("Name invalid, %w", err), errResp
	}

	err, errResp = isHardwareDescValid(req.Desc, true)
	if err != nil {
		return fmt.Errorf("Desc invalid, %w", err), errResp
	}

	err, errResp = isInstanceTypeValid(req.InstanceType, true)
	if err != nil {
		return fmt.Errorf("InstanceType invalid, %w", err), errResp
	}

	err, errResp = isInstanceFamilyValid(req.InstanceFamily, true)
	if err != nil {
		return fmt.Errorf("InstanceFamily invalid, %w", err), errResp
	}

	err, errResp = isNetworkValid(req.Network)
	if err != nil {
		return fmt.Errorf("Network invalid, %w", err), errResp
	}

	err, errResp = isCpuValid(req.Cpu)
	if err != nil {
		return fmt.Errorf("Cpu invalid, %w", err), errResp
	}

	err, errResp = isCpuModelValid(req.CpuModel, true)
	if err != nil {
		return fmt.Errorf("CpuModel invalid, %w", err), errResp
	}

	err, errResp = isMemValid(req.Mem)
	if err != nil {
		return fmt.Errorf("Mem invalid, %w", err), errResp
	}

	err, errResp = isGpuValid(req.Gpu)
	if err != nil {
		return fmt.Errorf("Gpu invalid, %w", err), errResp
	}

	err, errResp = isGpuModelValid(req.GpuModel, true)
	if err != nil {
		return fmt.Errorf("GpuModel invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func ValidateAdminPutHardwareRequest(req *hardware.AdminPutRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("put hardware request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isHardwareIdValid(req.HardwareId, false)
	if err != nil {
		return fmt.Errorf("HardwareId invalid, %w", err), errResp
	}

	err, errResp = isZoneValid(req.Zone, false)
	if err != nil {
		return fmt.Errorf("Zone invalid, %w", err), errResp
	}

	err, errResp = isHardwareNameValid(req.Name, false)
	if err != nil {
		return fmt.Errorf("Name invalid, %w", err), errResp
	}

	err, errResp = isHardwareDescValid(req.Desc, true)
	if err != nil {
		return fmt.Errorf("Desc invalid, %w", err), errResp
	}

	err, errResp = isInstanceTypeValid(req.InstanceType, false)
	if err != nil {
		return fmt.Errorf("InstanceType invalid, %w", err), errResp
	}

	err, errResp = isInstanceFamilyValid(req.InstanceFamily, true)
	if err != nil {
		return fmt.Errorf("InstanceFamily invalid, %w", err), errResp
	}

	err, errResp = isNetworkValid(req.Network)
	if err != nil {
		return fmt.Errorf("Network invalid, %w", err), errResp
	}

	err, errResp = isCpuValid(req.Cpu)
	if err != nil {
		return fmt.Errorf("Cpu invalid, %w", err), errResp
	}

	err, errResp = isCpuModelValid(req.CpuModel, true)
	if err != nil {
		return fmt.Errorf("CpuModel invalid, %w", err), errResp
	}

	err, errResp = isMemValid(req.Mem)
	if err != nil {
		return fmt.Errorf("Mem invalid, %w", err), errResp
	}

	err, errResp = isGpuValid(req.Gpu)
	if err != nil {
		return fmt.Errorf("Gpu invalid, %w", err), errResp
	}

	err, errResp = isGpuModelValid(req.GpuModel, true)
	if err != nil {
		return fmt.Errorf("GpuModel invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func ValidateAdminGetHardwareRequest(req *hardware.AdminGetRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("get hardware request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isHardwareIdValid(req.HardwareId, false)
	if err != nil {
		return fmt.Errorf("HardwareId invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func ValidateAdminListHardwareRequest(req *hardware.AdminListRequest, c *gin.Context) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("list hardware request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	// PageOffset PageSize 特殊处理一下，如果是查询参数 ?PageSize= ，gin的解析会将他解析为PageSize=0，而期望是赋默认值1000
	err, errResp = ensurePageOffset(&req.PageOffset, c)
	if err != nil {
		return fmt.Errorf("PageOffset invalid, %w", err), errResp
	}

	err, errResp = ensurePageSize(&req.PageSize, c)
	if err != nil {
		return fmt.Errorf("PageSize invalid, %w", err), errResp
	}

	err, errResp = isZoneValid(req.Zone, true)
	if err != nil {
		return fmt.Errorf("Zone invalid, %w", err), errResp
	}

	err, errResp = isHardwareNameValid(req.Name, true)
	if err != nil {
		return fmt.Errorf("Hardware Name invalid, %w", err), errResp
	}

	err, errResp = isCpuValid(req.Cpu)
	if err != nil {
		return fmt.Errorf("Cpu invalid, %w", err), errResp
	}

	err, errResp = isMemValid(req.Mem)
	if err != nil {
		return fmt.Errorf("Mem invalid, %w", err), errResp
	}

	err, errResp = isGpuValid(req.Gpu)
	if err != nil {
		return fmt.Errorf("Gpu invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func ValidateAdminDeleteHardwareRequest(req *hardware.AdminDeleteRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("delete hardware request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isHardwareIdValid(req.HardwareId, false)
	if err != nil {
		return fmt.Errorf("HardwareId invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func ValidateAdminPostHardwaresUsersRequest(req *hardware.AdminPostUsersRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("post hardwares users request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isHardwareIdsValid(req.Hardwares)
	if err != nil {
		return fmt.Errorf("Hardwares invalid, %w", err), errResp
	}

	err, errResp = isUserIdsValid(req.Users)
	if err != nil {
		return fmt.Errorf("Users invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func isHardwareIdsValid(hardwares []string) (error, response.ErrorResp) {
	if err := isSnowflakeIdsValid(hardwares); err != nil {
		return err, response.WrapErrorResp(common.InvalidArgumentHardwareId, err.Error())
	}

	return nil, response.ErrorResp{}
}

func ValidateAdminDeleteHardwaresUsersRequest(req *hardware.AdminDeleteUsersRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("delete hardwares users request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isHardwareIdsValid(req.Hardwares)
	if err != nil {
		return fmt.Errorf("Hardwares invalid, %w", err), errResp
	}

	err, errResp = isUserIdsValid(req.Users)
	if err != nil {
		return fmt.Errorf("Users invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}
func PInt(i int) *int {
	return &i
}
