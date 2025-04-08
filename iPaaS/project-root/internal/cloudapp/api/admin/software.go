package admin

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/rpc"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/utils"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/db"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/response"
)

func PostSoftwares(ctx *gin.Context) {
	logger := trace.GetLogger(ctx)

	req := new(softwareApi.AdminPostRequest)
	err := bindPostSoftwareRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind post software request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminPostSoftwareRequest(req)
	if err = response.BadRequestIfError(ctx, err, errResp); err != nil {
		logger.Warnf("validate admin post software request failed, %v", err)
		return
	}

	softwareId, err := rpc.GenID(ctx)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "generate snowflake id failed")); err != nil {
		logger.Warnf("generate snowflake id failed, %v", err)
		return
	}

	software := createSoftwareDBModelForPost(req, softwareId)
	err = dao.AddSoftware(ctx, software)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "add software to database failed")); err != nil {
		logger.Warnf("add software to database failed, %v", err)
		return
	}

	response.RenderJson(&softwareApi.AdminPostResponseData{
		SoftwareId: software.Id.String(),
	}, ctx)
}

func bindPostSoftwareRequest(req *softwareApi.AdminPostRequest, c *gin.Context) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("bind json failed, %w", err)
	}

	return nil
}

func createSoftwareDBModelForPost(req *softwareApi.AdminPostRequest, softwareId snowflake.ID) *models.Software {
	m := &models.Software{
		Id:         softwareId,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	if req == nil {
		return m
	}

	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Zone != nil {
		m.Zone = zone.Zone(*req.Zone)
	}
	if req.Desc != nil {
		m.Desc = *req.Desc
	}
	if req.Icon != nil {
		m.Icon = *req.Icon
	}
	if req.Platform != nil {
		m.Platform = models.Platform(*req.Platform)
	}
	if req.ImageId != nil {
		m.ImageId = *req.ImageId
	}
	if req.InitScript != nil {
		m.InitScript = *req.InitScript
	}
	if req.GpuDesired != nil {
		m.GpuDesired = req.GpuDesired
	} else {
		m.GpuDesired = utils.PBool(false)
	}

	return m
}

func PutSoftware(ctx *gin.Context) {
	logger := trace.GetLogger(ctx)

	req := new(softwareApi.AdminPutRequest)
	err := bindPutSoftwareRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind put software request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminPutSoftwareRequest(req)
	if err = response.BadRequestIfError(ctx, err, errResp); err != nil {
		logger.Warnf("validate admin put software request failed, %v", err)
		return
	}

	software := createSoftwareDBModelForPut(req)
	exist, err := dao.UpdateSoftwareAllCol(ctx, software)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "update software in database failed")); err != nil {
		logger.Warnf("update software in database failed, %v", err)
		return
	}
	if !exist {
		err = fmt.Errorf("software not found")
		_ = response.NotFoundIfError(ctx, err, response.WrapErrorResp(common.SoftwareNotFound, err.Error()))
		logger.Warn(err)
		return
	}

	response.RenderJson(nil, ctx)
}

func bindPutSoftwareRequest(req *softwareApi.AdminPutRequest, c *gin.Context) error {
	var err error
	if err = c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	if err = c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("bind json failed, %w", err)
	}

	return nil
}

func createSoftwareDBModelForPut(req *softwareApi.AdminPutRequest) *models.Software {
	m := &models.Software{
		UpdateTime: time.Now(),
	}
	if req == nil {
		return m
	}

	if req.SoftwareId != nil {
		m.Id = snowflake.MustParseString(*req.SoftwareId)
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Zone != nil {
		m.Zone = zone.Zone(*req.Zone)
	}
	if req.Desc != nil {
		m.Desc = *req.Desc
	}
	if req.Icon != nil {
		m.Icon = *req.Icon
	}
	if req.Platform != nil {
		m.Platform = models.Platform(*req.Platform)
	}
	if req.ImageId != nil {
		m.ImageId = *req.ImageId
	}
	if req.InitScript != nil {
		m.InitScript = *req.InitScript
	}
	if req.GpuDesired != nil {
		m.GpuDesired = req.GpuDesired
	} else {
		m.GpuDesired = utils.PBool(false)
	}

	return m
}

func PatchSoftware(ctx *gin.Context) {
	logger := trace.GetLogger(ctx)

	req := new(softwareApi.AdminPatchRequest)
	err := bindPatchSoftwareRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind patch software request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminPatchSoftwareRequest(req)
	if err = response.BadRequestIfError(ctx, err, errResp); err != nil {
		logger.Warnf("validate admin patch software request failed, %v", err)
		return
	}

	software, cols := generateSoftwareDBModelForPatch(req)
	exist, err := dao.UpdateSoftware(ctx, software, cols...)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "update software in database failed")); err != nil {
		logger.Warnf("update software in database failed, %v", err)
		return
	}
	if !exist {
		err = fmt.Errorf("software not found")
		_ = response.NotFoundIfError(ctx, err, response.WrapErrorResp(common.SoftwareNotFound, err.Error()))
		logger.Warn(err)
		return
	}

	response.RenderJson(nil, ctx)
}

func bindPatchSoftwareRequest(req *softwareApi.AdminPatchRequest, c *gin.Context) error {
	var err error
	if err = c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	if err = c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("bind json failed, %w", err)
	}

	return nil
}

func generateSoftwareDBModelForPatch(req *softwareApi.AdminPatchRequest) (*models.Software, []string) {
	software := &models.Software{
		Id: snowflake.MustParseString(*req.SoftwareId),
	}

	updateCols := make([]string, 0)

	if req.Zone != nil {
		software.Zone = zone.Zone(*req.Zone)
		updateCols = append(updateCols, "zone")
	}

	if req.Name != nil {
		software.Name = *req.Name
		updateCols = append(updateCols, "name")
	}

	if req.Desc != nil {
		software.Desc = *req.Desc
		updateCols = append(updateCols, "desc")
	}

	if req.Icon != nil {
		software.Icon = *req.Icon
		updateCols = append(updateCols, "icon")
	}

	if req.Platform != nil {
		software.Platform = models.Platform(*req.Platform)
		updateCols = append(updateCols, "platform")
	}

	if req.ImageId != nil {
		software.ImageId = *req.ImageId
		updateCols = append(updateCols, "image_id")
	}

	if req.InitScript != nil {
		software.InitScript = *req.InitScript
		updateCols = append(updateCols, "init_script")
	}

	if req.GpuDesired != nil {
		software.GpuDesired = req.GpuDesired
		updateCols = append(updateCols, "gpu_desired")
	}

	return software, updateCols
}

func GetSoftware(ctx *gin.Context) {
	logger := trace.GetLogger(ctx)

	req := new(softwareApi.AdminGetRequest)
	err := bindGetSoftwareRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind get software request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminGetSoftwareRequest(req)
	if err = response.BadRequestIfError(ctx, err, errResp); err != nil {
		logger.Warnf("validate admin get software request failed, %v", err)
		return
	}

	software, exist, err := dao.GetSoftware(ctx, snowflake.MustParseString(*req.SoftwareId))
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "get software from database failed")); err != nil {
		logger.Warnf("get software from datbase failed, %v", err)
		return
	}
	if !exist {
		err = fmt.Errorf("software not found")
		_ = response.NotFoundIfError(ctx, err, response.WrapErrorResp(common.SoftwareNotFound, err.Error()))
		logger.Warn(err)
		return
	}

	remoteApps, err := dao.ListRemoteAppBySoftwareID(ctx, software.Id)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "list remoteapp by softwareId failed")); err != nil {
		logger.Warnf("list remoteapp by softwareId failed, %v", err)
		return
	}

	remoteAppsResp := make([]*v20230530.RemoteApp, 0, len(remoteApps))
	for _, remoteApp := range remoteApps {
		remoteAppsResp = append(remoteAppsResp, remoteApp.ToHTTPModel())
	}

	data := &softwareApi.AdminGetResponseData{
		Software: v20230530.Software{
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
		},
	}

	response.RenderJson(data, ctx)
}

func bindGetSoftwareRequest(req *softwareApi.AdminGetRequest, c *gin.Context) error {
	if err := c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	return nil
}

func ListSoftware(ctx *gin.Context) {
	logger := trace.GetLogger(ctx)

	req := &softwareApi.AdminListRequest{}
	err := bindListSoftwareRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind list software request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminListSoftwareRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, errResp); err != nil {
		logger.Warnf("validate admin list software request failed, %v", err)
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

	params := ensureListSoftwareDaoParams(req)
	var softwares []*models.Software
	var total int64
	if isYSProduct || req.UserId == nil {
		// YSProduct用户不联表查，默认全能看到
		softwares, total, err = dao.ListSoftware(ctx, params)
	} else {
		softwares, total, err = dao.ListSoftwareByUser(ctx, params, userId)
	}
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "list software from db failed")); err != nil {
		logger.Warnf("list software from database failed, %v", err)
		return
	}

	data := &softwareApi.AdminListResponseData{
		Software: make([]*v20230530.Software, 0, len(softwares)),
		Offset:   params.PageOffset,
		Size:     params.PageSize,
		Total:    int(total),
	}

	for _, software := range softwares {
		remoteApps, err := dao.ListRemoteAppBySoftwareID(ctx, software.Id)
		if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "list remoteapp by softwareId failed")); err != nil {
			logger.Warnf("list remoteapp by softwareId failed, %v", err)
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

func bindListSoftwareRequest(req *softwareApi.AdminListRequest, c *gin.Context) error {
	if err := c.ShouldBindQuery(req); err != nil {
		return fmt.Errorf("bind query failed, %w", err)
	}

	return nil
}

func ensureListSoftwareDaoParams(req *softwareApi.AdminListRequest) *dao.ListSoftwareParams {
	params := &dao.ListSoftwareParams{
		PageOffset: *req.PageOffset,
		PageSize:   *req.PageSize,
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

	return params
}

func DeleteSoftware(c *gin.Context) {
	logger := trace.GetLogger(c)

	req := new(softwareApi.AdminDeleteRequest)
	err := bindDeleteSoftwareRequest(req, c)
	if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind delete software request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminDeleteSoftwareRequest(req)
	if err = response.BadRequestIfError(c, err, errResp); err != nil {
		logger.Warnf("validate admin delete software request failed, %v", err)
		return
	}

	exist := true
	softwareId := snowflake.MustParseString(*req.SoftwareId)
	err = with.DefaultTransaction(c, func(ctx context.Context) error {
		count, e := dao.DeleteSoftware(c, softwareId)
		if e != nil {
			return e
		}
		if count == 0 {
			exist = false
			return nil
		}

		// delete quota
		return dao.BatchDeleteSoftwareUsersBySoftwareId(ctx, []snowflake.ID{softwareId})
	})
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "delete software in database failed")); err != nil {
		logger.Warnf("delete software in database failed, %v", err)
		return
	}
	if !exist {
		err = fmt.Errorf("software not found")
		_ = response.NotFoundIfError(c, err, response.WrapErrorResp(common.SoftwareNotFound, err.Error()))
		logger.Warn(err)
		return
	}

	response.RenderJson(nil, c)
}

func bindDeleteSoftwareRequest(req *softwareApi.AdminDeleteRequest, c *gin.Context) error {
	if err := c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	return nil
}

func PostSoftwaresUsers(ctx *gin.Context) {
	logger := trace.GetLogger(ctx)

	req := new(softwareApi.AdminPostUsersRequest)
	err := bindPostSoftwaresUsersRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind post softwares users request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminPostSoftwaresUsersRequest(req)
	if err = response.BadRequestIfError(ctx, err, errResp); err != nil {
		logger.Warnf("validate admin post softwares users request, %v", err)
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

	softwares := util.MustParseToSnowflakeIds(req.Softwares)
	softwareExistList, err := dao.BatchCheckSoftwaresExist(ctx, softwares)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "batch check softwares exist failed")); err != nil {
		logger.Warnf("batch check softwares exist failed, %v", err)
		return
	}

	softwareNonexistentList := util.GetNonexistentList(softwares, softwareExistList)
	if len(softwareNonexistentList) > 0 {
		err = fmt.Errorf("softwares %v not exist", softwareNonexistentList)
		_ = response.NotFoundIfError(ctx, err, response.WrapErrorResp(common.SoftwareNotFound, err.Error()))
		logger.Warn(err)
		return
	}

	err = dao.BatchAddSoftwareUsers(ctx, softwares, users)
	if err != nil {
		if errors.Is(err, db.ErrDuplicatedEntry) {
			// 查一下是哪些已经存在了
			softwareUsers, e := dao.GetSofwareUserByUsers(ctx, users)
			if e = response.InternalErrorIfError(ctx, e, response.WrapErrorResp(common.InternalServerErrorCode, "get software users failed")); e != nil {
				return
			}

			// 与入参做交集，抛错指出已存在部分
			errMsg := softwareUsersDuplicatedErrMsg(intersectSoftwareUsers(req, softwareUsers))
			_ = response.ConflictIfError(ctx, err, response.WrapErrorResp(common.InvalidArgumentErrorCode, errMsg))
			logger.Warnf("batch add software users to database failed by duplication, %v, errMsg: %s", err, errMsg)
			return
		}

		_ = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "batch add softwares users to database failed"))
		logger.Warnf("batch add softwares users to database failed, %v", err)
		return
	}

	response.RenderJson(nil, ctx)
}

func bindPostSoftwaresUsersRequest(req *softwareApi.AdminPostUsersRequest, c *gin.Context) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("bind json failed, %w", err)
	}

	return nil
}

func intersectSoftwareUsers(req *softwareApi.AdminPostUsersRequest, existInDB []models.SoftwareUser) []models.SoftwareUser {
	setInReq := make(map[snowflake.ID]snowflake.ID)
	for _, software := range req.Softwares {
		for _, user := range req.Users {
			setInReq[snowflake.MustParseString(software)] = snowflake.MustParseString(user)
		}
	}

	res := make([]models.SoftwareUser, 0)
	for software, user := range setInReq {
		for _, v := range existInDB {
			if software == v.SoftwareId && user == v.UserId {
				res = append(res, models.SoftwareUser{
					SoftwareId: software,
					UserId:     user,
				})
			}
		}
	}

	return res
}

func softwareUsersDuplicatedErrMsg(softwareUsersDuplicated []models.SoftwareUser) string {
	errMsg := "those record already exist, [softwareId, userId]: "
	for _, v := range softwareUsersDuplicated {
		errMsg += fmt.Sprintf(" [%s, %s] ", v.SoftwareId, v.UserId)
	}

	return errMsg
}

func DeleteSoftwaresUsers(ctx *gin.Context) {
	logger := trace.GetLogger(ctx)

	req := new(softwareApi.AdminDeleteUsersRequest)
	err := bindDeleteSoftwaresUsersRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind delete softwares users request, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminDeleteSoftwaresUsersRequest(req)
	if err = response.BadRequestIfError(ctx, err, errResp); err != nil {
		logger.Warnf("validate admin delete softwares users request failed, %v", err)
		return
	}

	users := util.MustParseToSnowflakeIds(req.Users)
	softwares := util.MustParseToSnowflakeIds(req.Softwares)
	err = dao.BatchDeleteSoftwareUsers(ctx, softwares, users)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "batch delete softwares users in database failed")); err != nil {
		logger.Warnf("batch delete softwares users to database failed, %v", err)
		return
	}

	response.RenderJson(nil, ctx)
}

func bindDeleteSoftwaresUsersRequest(req *softwareApi.AdminDeleteUsersRequest, c *gin.Context) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("bind json failed, %w", err)
	}

	return nil
}
