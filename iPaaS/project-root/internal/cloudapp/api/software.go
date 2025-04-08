package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"

	softwareApi "github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/software"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/validator"
	zone "github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/response"
)

func GetSoftWare(ctx *gin.Context) {
	logger := trace.GetLogger(ctx).Base()

	userId, err := util.GetUserId(ctx)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.UserNotExistsErrorCode, "invalid userId")); err != nil {
		logger.Warnf("get userId failed, %v", err)
		return
	}
	logger = logger.With("user-id", userId.String())

	req := new(softwareApi.APIGetRequest)
	err = bindGetSoftwareRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind get software request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAPIGetSoftwareRequest(req)
	if err = response.BadRequestIfError(ctx, err, errResp); err != nil {
		logger.Warnf("validate api get software request failed, %v", err)
		return
	}

	state, err := util.GetState(ctx)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "get state failed")); err != nil {
		logger.Warnf("get state from gin ctx failed, %v", err)
		return
	}

	isYSProduct, err := state.IamClient.IsYsProductUser(userId)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "check user is YSProduct or not failed")); err != nil {
		logger.Warnf("check user is YSProductUser by iam client failed, %v", err)
		return
	}

	softwareId := snowflake.MustParseString(*req.SoftwareId)
	var software *models.Software
	var exist bool
	if isYSProduct {
		// YSProduct用户不联表查，默认全能看到
		software, exist, err = dao.GetSoftware(ctx, softwareId)
	} else {
		software, exist, err = dao.GetSoftwareByUser(ctx, softwareId, userId)
	}
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "get software from database failed")); err != nil {
		logger.Warnf("get software from database failed, %v", err)
		return
	}
	if !exist {
		err = fmt.Errorf("software not found")
		_ = response.NotFoundIfError(ctx, err, response.WrapErrorResp(common.SoftwareNotFound, err.Error()))
		logger.Warn(err)
		return
	}

	remoteApps, err := dao.ListRemoteAppBySoftwareID(ctx, software.Id)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "database error")); err != nil {
		logger.Error("list remote app by software id failed, %v", err)
		return
	}

	remoteAppsResp := make([]*v20230530.RemoteApp, 0, len(remoteApps))
	for _, remoteApp := range remoteApps {
		remoteAppsResp = append(remoteAppsResp, remoteApp.ToHTTPModel())
	}

	softwareData := *software.ToHTTPModel()
	softwareData.RemoteApps = remoteAppsResp

	data := &softwareApi.APIGetResponseData{
		Software: softwareData,
	}

	response.RenderJson(data, ctx)
}

func bindGetSoftwareRequest(req *softwareApi.APIGetRequest, c *gin.Context) error {
	if err := c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	return nil
}

func ListSoftWare(ctx *gin.Context) {
	logger := trace.GetLogger(ctx).Base()

	userId, err := util.GetUserId(ctx)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.UserNotExistsErrorCode, "invalid userId")); err != nil {
		logger.Warnf("get userId failed, %v", err)
		return
	}
	logger = logger.With("user-id", userId.String())

	req := new(softwareApi.APIListRequest)
	err = bindListSoftwareRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind list software request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAPIListSoftwareRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, errResp); err != nil {
		logger.Warnf("validate api list software request failed, %v", err)
		return
	}

	state, err := util.GetState(ctx)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "get state failed")); err != nil {
		logger.Warnf("get state from gin ctx failed, %v", err)
		return
	}

	isYSProduct, err := state.IamClient.IsYsProductUser(userId)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "check user is YSProduct or not failed")); err != nil {
		logger.Warnf("check user is YSProductUser by iam client failed, %v", err)
		return
	}

	params := ensureListSoftwareDaoParams(req)
	var softwares []*models.Software
	var total int64
	if isYSProduct {
		// YSProduct用户不联表查，默认全能看到
		softwares, total, err = dao.ListSoftware(ctx, params)
	} else {
		softwares, total, err = dao.ListSoftwareByUser(ctx, params, userId)
	}
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "database error")); err != nil {
		logger.Warnf("list software failed, %v", err)
		return
	}

	data := &softwareApi.APIListResponseData{
		Software: make([]*v20230530.Software, 0, len(softwares)),
		Offset:   params.PageOffset,
		Size:     params.PageSize,
		Total:    int(total),
	}

	for _, software := range softwares {
		remoteApps, err := dao.ListRemoteAppBySoftwareID(ctx, software.Id)
		if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "database error")); err != nil {
			logger.Warnf("list remote app by software id [%s] failed, %v", software.Id, err)
			return
		}

		remoteAppsResp := make([]*v20230530.RemoteApp, 0, len(remoteApps))
		for _, remoteApp := range remoteApps {
			remoteAppsResp = append(remoteAppsResp, remoteApp.ToHTTPModel())
		}

		data.Software = append(data.Software, &v20230530.Software{
			SoftwareId: software.Id.String(),
			Zone:       software.Zone.String(),
			Name:       software.Name,
			Desc:       software.Desc,
			Icon:       software.Icon,
			Platform:   string(software.Platform),
			ImageId:    software.ImageId,
			InitScript: software.InitScript,
			GpuDesired: software.GpuDesired,
			RemoteApps: remoteAppsResp,
		})
	}

	if params.PageSize+params.PageOffset < int(total) {
		data.NextMarker = params.PageOffset + params.PageSize
	} else {
		data.NextMarker = -1
	}

	response.RenderJson(data, ctx)
}

func bindListSoftwareRequest(req *softwareApi.APIListRequest, c *gin.Context) error {
	if err := c.ShouldBindQuery(req); err != nil {
		return fmt.Errorf("bind query failed, %w", err)
	}

	return nil
}

func ensureListSoftwareDaoParams(req *softwareApi.APIListRequest) *dao.ListSoftwareParams {
	params := &dao.ListSoftwareParams{}
	if req == nil {
		return params
	}

	if req.Name != nil {
		params.Name = *req.Name
	}

	if req.Platform != nil {
		params.Platform = *req.Platform
	}

	if req.Zone != nil {
		params.Zone = zone.Zone(*req.Zone)
	}

	if req.PageOffset != nil {
		params.PageOffset = *req.PageOffset
	}

	if req.PageSize != nil {
		params.PageSize = *req.PageSize
	}

	return params
}
