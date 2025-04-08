package admin

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"
	hardwareApi "github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/hardware"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/validator"
	zone "github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/rpc"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/db"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/response"
)

func PostHardwares(ctx *gin.Context) {
	logger := trace.GetLogger(ctx)

	req := new(hardwareApi.AdminPostRequest)
	err := bindPostHardwareRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind post hardware request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminPostHardwaresRequest(req)
	if err = response.BadRequestIfError(ctx, err, errResp); err != nil {
		logger.Warnf("validate admin post hardwares request failed, %v", err)
		return
	}

	hardwareId, err := rpc.GenID(ctx)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "generate snowflake id failed")); err != nil {
		logger.Warnf("generate snowflake id failed, %v", err)
		return
	}

	hardware := createHardwareDBModelForPost(req, hardwareId)
	err = dao.AddHardware(ctx, hardware)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "add hardware to database failed")); err != nil {
		logger.Warnf("add hardware to database failed, %v", err)
		return
	}

	response.RenderJson(hardwareApi.AdminPostResponseData{
		HardwareId: hardware.Id.String(),
	}, ctx)
}

func createHardwareDBModelForPost(req *hardwareApi.AdminPostRequest, hardwareId snowflake.ID) *models.Hardware {
	m := &models.Hardware{
		Id:         hardwareId,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	if req == nil {
		return m
	}

	if req.Zone != nil {
		m.Zone = zone.Zone(*req.Zone)
	}

	if req.Name != nil {
		m.Name = *req.Name
	}

	if req.Desc != nil {
		m.Desc = *req.Desc
	}

	if req.InstanceType != nil {
		m.InstanceType = *req.InstanceType
	}

	if req.InstanceFamily != nil {
		m.InstanceFamily = *req.InstanceFamily
	}

	if req.Network != nil {
		m.Network = int64(*req.Network)
	}

	if req.Cpu != nil {
		m.Cpu = int64(*req.Cpu)
	}

	if req.CpuModel != nil {
		m.CpuModel = *req.CpuModel
	}

	if req.Mem != nil {
		m.Mem = int64(*req.Mem)
	}

	if req.Gpu != nil {
		m.Gpu = int64(*req.Gpu)
	}

	if req.GpuModel != nil {
		m.GpuModel = *req.GpuModel
	}

	return m
}

func bindPostHardwareRequest(req *hardwareApi.AdminPostRequest, c *gin.Context) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("bind json failed, %w", err)
	}

	return nil
}

func PatchHardware(ctx *gin.Context) {
	logger := trace.GetLogger(ctx)

	req := new(hardwareApi.AdminPatchRequest)
	err := bindPatchHardwareRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind patch hardware request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminPatchHardwareRequest(req)
	if err = response.BadRequestIfError(ctx, err, errResp); err != nil {
		logger.Warnf("validate admin patch hardware request failed, %v", err)
		return
	}

	hardwareId := snowflake.MustParseString(*req.HardwareId)
	hardware, cols := generateHardwareDBModelForPatch(req)
	err = dao.UpdateHardware(ctx, hardwareId, hardware, cols...)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "update hardware to database failed")); err != nil {
		logger.Warnf("update hardware failed, %v", err)
		return
	}

	response.RenderJson(nil, ctx)
}

func bindPatchHardwareRequest(req *hardwareApi.AdminPatchRequest, c *gin.Context) error {
	var err error
	if err = c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	if err = c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("bind json failed, %w", err)
	}

	return nil
}

func generateHardwareDBModelForPatch(req *hardwareApi.AdminPatchRequest) (*models.Hardware, []string) {
	hardware := &models.Hardware{}

	updateCols := make([]string, 0)

	if req.Zone != nil {
		hardware.Zone = zone.Zone(*req.Zone)
		updateCols = append(updateCols, "zone")
	}

	if req.Name != nil {
		hardware.Name = *req.Name
		updateCols = append(updateCols, "name")
	}

	if req.Desc != nil {
		hardware.Desc = *req.Desc
		updateCols = append(updateCols, "desc")
	}

	if req.InstanceType != nil {
		hardware.InstanceType = *req.InstanceType
		updateCols = append(updateCols, "instance_type")
	}

	if req.InstanceFamily != nil {
		hardware.InstanceFamily = *req.InstanceFamily
		updateCols = append(updateCols, "instance_family")
	}

	if req.Network != nil {
		hardware.Network = int64(*req.Network)
		updateCols = append(updateCols, "network")
	}

	if req.Cpu != nil {
		hardware.Cpu = int64(*req.Cpu)
		updateCols = append(updateCols, "cpu")
	}

	if req.Mem != nil {
		hardware.Mem = int64(*req.Mem)
		updateCols = append(updateCols, "mem")
	}

	if req.Gpu != nil {
		hardware.Gpu = int64(*req.Gpu)
		updateCols = append(updateCols, "gpu")
	}

	if req.GpuModel != nil {
		hardware.GpuModel = *req.GpuModel
		updateCols = append(updateCols, "gpu_model")
	}

	if req.CpuModel != nil {
		hardware.CpuModel = *req.CpuModel
		updateCols = append(updateCols, "cpu_model")
	}

	return hardware, updateCols
}

func PutHardware(ctx *gin.Context) {
	logger := trace.GetLogger(ctx)

	req := new(hardwareApi.AdminPutRequest)
	err := bindPutHardwareRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind put hardware request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminPutHardwareRequest(req)
	if err = response.BadRequestIfError(ctx, err, errResp); err != nil {
		logger.Warnf("validate admin put hardware request failed, %v", err)
		return
	}

	hardwareId := snowflake.MustParseString(*req.HardwareId)
	hardware := createHardwareDBModelForPut(req)

	err = dao.UpdateHardwareAllCol(ctx, hardwareId, hardware)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "update hardware in database failed")); err != nil {
		logger.Warnf("update hardware in database failed, %v", err)
		return
	}

	response.RenderJson(nil, ctx)
}

func bindPutHardwareRequest(req *hardwareApi.AdminPutRequest, c *gin.Context) error {
	var err error
	if err = c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	if err = c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("bind json failed, %w", err)
	}

	return nil
}

func createHardwareDBModelForPut(req *hardwareApi.AdminPutRequest) *models.Hardware {
	m := &models.Hardware{
		UpdateTime: time.Now(),
	}
	if req == nil {
		return m
	}

	if req.Zone != nil {
		m.Zone = zone.Zone(*req.Zone)
	}

	if req.Name != nil {
		m.Name = *req.Name
	}

	if req.Desc != nil {
		m.Desc = *req.Desc
	}

	if req.InstanceType != nil {
		m.InstanceType = *req.InstanceType
	}

	if req.InstanceFamily != nil {
		m.InstanceFamily = *req.InstanceFamily
	}

	if req.Network != nil {
		m.Network = int64(*req.Network)
	}

	if req.Cpu != nil {
		m.Cpu = int64(*req.Cpu)
	}

	if req.CpuModel != nil {
		m.CpuModel = *req.CpuModel
	}

	if req.Mem != nil {
		m.Mem = int64(*req.Mem)
	}

	if req.Gpu != nil {
		m.Gpu = int64(*req.Gpu)
	}

	if req.GpuModel != nil {
		m.GpuModel = *req.GpuModel
	}

	return m
}

func GetHardware(ctx *gin.Context) {
	logger := trace.GetLogger(ctx)

	req := new(hardwareApi.AdminGetRequest)
	err := bindGetHardwareRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind get hardware request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminGetHardwareRequest(req)
	if err = response.BadRequestIfError(ctx, err, errResp); err != nil {
		logger.Warnf("validate admin get hardware request failed, %v", err)
		return
	}

	hardewareId := snowflake.MustParseString(*req.HardwareId)
	hardware, exist, err := dao.GetHardware(ctx, hardewareId)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "get hardware from database failed")); err != nil {
		logger.Warnf("get hardware from database failed, %v", err)
		return
	}
	if !exist {
		err = fmt.Errorf("hardware not found")
		_ = response.NotFoundIfError(ctx, err, response.WrapErrorResp(common.HardwareNotFound, err.Error()))
		logger.Warn(err)
		return
	}

	data := &hardwareApi.AdminGetResponseData{
		Hardware: v20230530.Hardware{
			HardwareId:     hardware.Id.String(),
			Zone:           hardware.Zone.String(),
			Name:           hardware.Name,
			Desc:           hardware.Desc,
			InstanceType:   hardware.InstanceType,
			InstanceFamily: hardware.InstanceFamily,
			Network:        int(hardware.Network),
			Cpu:            int(hardware.Cpu),
			Mem:            int(hardware.Mem),
			Gpu:            int(hardware.Gpu),
			GpuModel:       hardware.GpuModel,
			CpuModel:       hardware.CpuModel,
		},
	}

	response.RenderJson(data, ctx)
}

func bindGetHardwareRequest(req *hardwareApi.AdminGetRequest, c *gin.Context) error {
	if err := c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	return nil
}

func ListHardware(ctx *gin.Context) {
	logger := trace.GetLogger(ctx)

	req := new(hardwareApi.AdminListRequest)
	err := bindListHardwareRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind list hardware request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminListHardwareRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, errResp); err != nil {
		logger.Warnf("validate admin list hardware request failed, %v", err)
		return
	}

	var userId snowflake.ID
	var isYSProduct bool
	if req.UserId != nil {

		state, err := util.GetState(ctx)
		if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "get state failed")); err != nil {
			logger.Warnf("get state from gin ctx failed, %v", err)
			return
		}
		if userId, err = snowflake.ParseString(*req.UserId); err != nil {
			logger.Warnf("parse userId %s to snowflake id failed, %v", *req.UserId, err)
			return
		}

		isYSProduct, err = state.IamClient.IsYsProductUser(userId)
		if err != nil {
			if strings.Contains(err.Error(), "user not found") {
				logger.Infof("check user is YSProduct User by iam client failed, user not found")
			} else {
				logger.Warnf("check user is YSProduct User by iam client failed, %v", err)
			}
		}
	}

	params := ensureListHardwareParams(req)
	var hardwares []*models.Hardware
	var total int64
	if isYSProduct || req.UserId == nil {
		// YSProduct用户不联表查，默认全能看到
		hardwares, total, err = dao.ListHardware(ctx, params)
	} else {
		hardwares, total, err = dao.ListHardwareByUser(ctx, params, userId)
	}
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "list hardware in database failed")); err != nil {
		logger.Warnf("list hardware in database failed, %v", err)
		return
	}

	data := &hardwareApi.AdminListResponseData{
		Hardware: make([]*v20230530.Hardware, 0, len(hardwares)),
		Offset:   params.PageOffset,
		Size:     params.PageSize,
		Total:    int(total),
	}

	if params.PageSize+params.PageOffset < int(total) {
		data.NextMarker = params.PageOffset + params.PageSize
	} else {
		data.NextMarker = -1
	}

	for _, hardware := range hardwares {
		data.Hardware = append(data.Hardware, &v20230530.Hardware{
			HardwareId:     hardware.Id.String(),
			Zone:           hardware.Zone.String(),
			Name:           hardware.Name,
			Desc:           hardware.Desc,
			InstanceType:   hardware.InstanceType,
			InstanceFamily: hardware.InstanceFamily,
			Network:        int(hardware.Network),
			Cpu:            int(hardware.Cpu),
			Mem:            int(hardware.Mem),
			Gpu:            int(hardware.Gpu),
			GpuModel:       hardware.GpuModel,
			CpuModel:       hardware.CpuModel,
		})
	}

	response.RenderJson(data, ctx)
}

func bindListHardwareRequest(req *hardwareApi.AdminListRequest, c *gin.Context) error {
	if err := c.ShouldBindQuery(req); err != nil {
		return fmt.Errorf("bind query failed, %w", err)
	}

	return nil
}

func ensureListHardwareParams(req *hardwareApi.AdminListRequest) *dao.ListHardwareParams {
	params := &dao.ListHardwareParams{}
	if req == nil {
		return params
	}

	if req.Name != nil {
		params.Name = *req.Name
	}

	if req.Zone != nil {
		params.Zone = zone.Zone(*req.Zone)
	}

	if req.Cpu != nil {
		params.Cpu = *req.Cpu
	}

	if req.Mem != nil {
		params.Mem = *req.Mem
	}

	if req.Gpu != nil {
		params.Gpu = *req.Gpu
	}

	if req.PageOffset != nil {
		params.PageOffset = *req.PageOffset
	}

	if req.PageSize != nil {
		params.PageSize = *req.PageSize
	}

	return params
}

func DeleteHardware(c *gin.Context) {
	logger := trace.GetLogger(c)

	req := new(hardwareApi.AdminDeleteRequest)
	err := bindDeleteHardwareRequest(req, c)
	if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind delete hardware request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminDeleteHardwareRequest(req)
	if err = response.BadRequestIfError(c, err, errResp); err != nil {
		logger.Warnf("validate admin delete hardware request failed, %v", err)
		return
	}

	hardwareId := snowflake.MustParseString(*req.HardwareId)
	exist := true
	err = with.DefaultTransaction(c, func(ctx context.Context) error {
		count, e := dao.DeleteHardware(c, hardwareId)
		if e != nil {
			return e
		}
		if count == 0 {
			exist = false
			return nil
		}

		// delete quota
		return dao.BatchDeleteHardwareUsersByHardwareId(c, []snowflake.ID{hardwareId})
	})
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "delete hardware from database failed")); err != nil {
		logger.Warnf("delete hardware from database failed, %v", err)
		return
	}
	if !exist {
		err = fmt.Errorf("hardware not found")
		_ = response.NotFoundIfError(c, err, response.WrapErrorResp(common.HardwareNotFound, err.Error()))
		logger.Warn(err)
		return
	}

	response.RenderJson(nil, c)
}

func bindDeleteHardwareRequest(req *hardwareApi.AdminDeleteRequest, c *gin.Context) error {
	if err := c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	return nil
}

func PostHardwaresUsers(ctx *gin.Context) {
	logger := trace.GetLogger(ctx)

	req := new(hardwareApi.AdminPostUsersRequest)
	err := bindPostHardwaresUsersRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind post hardwares users request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminPostHardwaresUsersRequest(req)
	if err = response.BadRequestIfError(ctx, err, errResp); err != nil {
		logger.Warnf("validate admin post hardwares request failed, %v", err)
		return
	}

	users := util.MustParseToSnowflakeIds(req.Users)
	userExistList, err := rpc.BatchCheckUserExist(ctx, req.Users)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "call hydra-lcp failed")); err != nil {
		logger.Warnf("call hydra-lcp to batch check user exist failed, %v", err)
		return
	}

	userNonExistentList := util.GetNonexistentList(users, userExistList)
	if len(userNonExistentList) > 0 {
		err = fmt.Errorf("users %v not exist", userNonExistentList)
		_ = response.NotFoundIfError(ctx, err, response.WrapErrorResp(common.UserNotExistsErrorCode, err.Error()))
		logger.Warn(err)
		return
	}

	hardwares := util.MustParseToSnowflakeIds(req.Hardwares)
	hardwareExistList, err := dao.BatchCheckHardwaresExist(ctx, hardwares)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "check hardware exist failed")); err != nil {
		logger.Warnf("batch check hardwares exist failed, %v", err)
		return
	}

	hardwareNonexistentList := util.GetNonexistentList(hardwares, hardwareExistList)
	if len(hardwareNonexistentList) > 0 {
		err = fmt.Errorf("hardwares %v not exist", hardwareNonexistentList)
		_ = response.NotFoundIfError(ctx, err, response.WrapErrorResp(common.HardwareNotFound, err.Error()))
		logger.Warn(err)
		return
	}

	err = dao.BatchAddHardwareUsers(ctx, hardwares, users)
	if err != nil {
		if errors.Is(err, db.ErrDuplicatedEntry) {
			// 查一下是哪些已经存在了
			hardwareUsers, e := dao.GetHardwareUserByUsers(ctx, users)
			if e = response.InternalErrorIfError(ctx, e, response.WrapErrorResp(common.InternalServerErrorCode, "get hardware users failed")); e != nil {
				return
			}

			// 与入参做交集，跑错指出已存在部分
			errMsg := hardwareUsersDuplicatedErrMsg(intersectHardwareUsers(req, hardwareUsers))
			_ = response.ConflictIfError(ctx, err, response.WrapErrorResp(common.InvalidArgumentErrorCode, errMsg))
			logger.Warnf("batch add hardware users to database failed by duplication, %v, errMsg: %s", err, errMsg)
			return
		}

		_ = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "batch add hardware users to database failed"))
		logger.Warnf("batch add hardware users to database failed, %v", err)
		return
	}

	response.RenderJson(nil, ctx)
}

func bindPostHardwaresUsersRequest(req *hardwareApi.AdminPostUsersRequest, c *gin.Context) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("bind json failed, %w", err)
	}

	return nil
}

func intersectHardwareUsers(req *hardwareApi.AdminPostUsersRequest, existInDB []models.HardwareUser) []models.HardwareUser {
	setInReq := make(map[snowflake.ID]snowflake.ID)
	for _, hardware := range req.Hardwares {
		for _, user := range req.Users {
			setInReq[snowflake.MustParseString(hardware)] = snowflake.MustParseString(user)
		}
	}

	res := make([]models.HardwareUser, 0)
	for hardware, user := range setInReq {
		for _, v := range existInDB {
			if hardware == v.HardwareId && user == v.UserId {
				res = append(res, models.HardwareUser{
					HardwareId: hardware,
					UserId:     user,
				})
			}
		}
	}

	return res
}

func hardwareUsersDuplicatedErrMsg(hardwareUsersDuplicated []models.HardwareUser) string {
	errMsg := "those record already exist, [hardwareId, userId]: "
	for _, v := range hardwareUsersDuplicated {
		errMsg += fmt.Sprintf(" [%s, %s] ", v.HardwareId, v.UserId)
	}

	return errMsg
}

func DeleteHardwaresUsers(ctx *gin.Context) {
	logger := trace.GetLogger(ctx)

	req := new(hardwareApi.AdminDeleteUsersRequest)
	err := bindDeleteHardwaresUsersRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind delete hardwares users request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminDeleteHardwaresUsersRequest(req)
	if err = response.BadRequestIfError(ctx, err, errResp); err != nil {
		logger.Warnf("validate admin delete hardwares users request failed, %v", err)
		return
	}

	users := util.MustParseToSnowflakeIds(req.Users)
	hardwares := util.MustParseToSnowflakeIds(req.Hardwares)
	err = dao.BatchDeleteHardwareUsers(ctx, hardwares, users)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "batch delete hardware users in database failed")); err != nil {
		logger.Warnf("batch delete hardware users in database failed, %v", err)
		return
	}
}

func bindDeleteHardwaresUsersRequest(req *hardwareApi.AdminDeleteUsersRequest, c *gin.Context) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("bind json failed, %w", err)
	}

	return nil
}
