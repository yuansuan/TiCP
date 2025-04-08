package admin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"
	remoteappmodel "github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/remoteapp"
	"github.com/yuansuan/ticp/common/project-root-api/common"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/rpc"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/db"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/response"
)

func PostRemoteApps(c *gin.Context) {
	logger := trace.GetLogger(c)

	req := new(remoteappmodel.AdminPostRequest)
	err := bindPostRemoteAppsRequest(req, c)
	if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind json failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminPostRemoteAppsRequest(req)
	if err = response.BadRequestIfError(c, err, errResp); err != nil {
		logger.Warnf("validate admin post remote apps request failed, %v", err)
		return
	}

	remoteApp, err, errResp := generateRemoteAppDBModelForPost(req, c)
	if err = response.InternalErrorIfError(c, err, errResp); err != nil {
		logger.Warnf("generate remote app db model failed, %v", err)
		return
	}

	err = dao.AddRemoteApp(c, remoteApp)
	if err != nil {
		if errors.Is(err, db.ErrDuplicatedEntry) {
			_ = response.BadRequestIfError(c, err, response.WrapErrorResp(common.InvalidArgumentName, fmt.Sprintf("remoteapp name [%s] already exist", remoteApp.Name)))
			logger.Warnf("remoteapp name [%s] already exist", remoteApp.Name)
			return
		}

		_ = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "add remote app to db failed"))
		logger.Warnf("add remote app to db failed, %v", err)
		return
	}

	response.RenderJson(remoteappmodel.AdminPostResponseData{
		Id: remoteApp.Id.String(),
	}, c)
}

func bindPostRemoteAppsRequest(req *remoteappmodel.AdminPostRequest, c *gin.Context) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("bind json failed, %w", err)
	}

	return nil
}

func generateRemoteAppDBModelForPost(req *remoteappmodel.AdminPostRequest, c *gin.Context) (*models.RemoteApp, error, response.ErrorResp) {
	id, err := rpc.GenID(c)
	if err != nil {
		return nil, fmt.Errorf("genereate snowflake id failed, %w", err), response.WrapErrorResp(common.InternalServerErrorCode, "generate snowflake id failed")
	}

	remoteApp := &models.RemoteApp{
		Id:         id,
		SoftwareId: snowflake.MustParseString(*req.SoftwareId),
	}

	if req.Desc != nil {
		remoteApp.Desc = *req.Desc
	}
	if req.Name != nil {
		remoteApp.Name = *req.Name
	}
	if req.Dir != nil {
		remoteApp.Dir = *req.Dir
	}
	if req.Args != nil {
		remoteApp.Args = *req.Args
	}
	if req.Logo != nil {
		remoteApp.Logo = *req.Logo
	}
	if req.DisableGfx != nil {
		remoteApp.DisableGfx = *req.DisableGfx
	}
	if req.LoginUser != nil {
		remoteApp.LoginUser = *req.LoginUser
	}

	return remoteApp, nil, response.ErrorResp{}
}

var remoteAppAllColumn = []string{"software_id", "desc", "base_url", "name", "dir", "args", "logo", "disable_gfx", "login_user"}

func PutRemoteApp(c *gin.Context) {
	logger := trace.GetLogger(c)

	req := new(remoteappmodel.AdminPutRequest)
	err := bindPutRemoteAppRequest(req, c)
	if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind json failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminPutRemoteAppRequest(req)
	if err = response.BadRequestIfError(c, err, errResp); err != nil {
		logger.Warnf("validate admin put remote app request failed, %v", err)
		return
	}

	remoteApp, err, errResp := generateRemoteAppDBModelForPut(snowflake.MustParseString(*req.RemoteAppId), req)
	if err = response.BadRequestIfError(c, err, errResp); err != nil {
		logger.Warnf("generate remote app db model for put failed, %v", err)
		return
	}

	exist, err := dao.UpdateRemoteApp(c, remoteApp, remoteAppAllColumn...)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "update remote app failed")); err != nil {
		logger.Warnf("update remote app failed, %v", err)
		return
	}
	if !exist {
		err = fmt.Errorf("remoteapp not found")
		_ = response.NotFoundIfError(c, err, response.WrapErrorResp(common.RemoteAppNotFound, err.Error()))
		logger.Warn(err)
		return
	}

	response.RenderJson(nil, c)
}

func bindPutRemoteAppRequest(req *remoteappmodel.AdminPutRequest, c *gin.Context) error {
	var err error
	if err = c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	if err = c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("bind json failed, %w", err)
	}

	return nil
}

func generateRemoteAppDBModelForPut(remoteAppId snowflake.ID, req *remoteappmodel.AdminPutRequest) (*models.RemoteApp, error, response.ErrorResp) {
	remoteApp := &models.RemoteApp{
		Id:         remoteAppId,
		SoftwareId: snowflake.MustParseString(*req.SoftwareId),
	}

	if req.Desc != nil {
		remoteApp.Desc = *req.Desc
	}
	if req.Name != nil {
		remoteApp.Name = *req.Name
	}
	if req.Dir != nil {
		remoteApp.Dir = *req.Dir
	}
	if req.Args != nil {
		remoteApp.Args = *req.Args
	}
	if req.Logo != nil {
		remoteApp.Logo = *req.Logo
	}
	if req.DisableGfx != nil {
		remoteApp.DisableGfx = *req.DisableGfx
	}
	if req.LoginUser != nil {
		remoteApp.LoginUser = *req.LoginUser
	}

	return remoteApp, nil, response.ErrorResp{}
}

func PatchRemoteApp(c *gin.Context) {
	logger := trace.GetLogger(c)

	req := new(remoteappmodel.AdminPatchRequest)
	err := bingPatchRemoteAppRequest(req, c)
	if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind json failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminPatchRemoteAppRequest(req)
	if err = response.BadRequestIfError(c, err, errResp); err != nil {
		logger.Warnf("validate admin patch remote app request failed, %v", err)
		return
	}

	remoteApp, updateCols := generateRemoteAppDBModelForPatch(snowflake.MustParseString(*req.RemoteAppId), req)
	exist, err := dao.UpdateRemoteApp(c, remoteApp, updateCols...)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "update remote app failed")); err != nil {
		logger.Warnf("update remote app failed, %v", err)
		return
	}
	if !exist {
		err = fmt.Errorf("remoteapp not found")
		_ = response.NotFoundIfError(c, err, response.WrapErrorResp(common.RemoteAppNotFound, err.Error()))
		logger.Warn(err)
		return
	}

	response.RenderJson(nil, c)
}

func bingPatchRemoteAppRequest(req *remoteappmodel.AdminPatchRequest, c *gin.Context) error {
	var err error
	if err = c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	if err = c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("bind json failed, %w", err)
	}

	return nil
}

func generateRemoteAppDBModelForPatch(remoteAppId snowflake.ID, req *remoteappmodel.AdminPatchRequest) (*models.RemoteApp, []string) {
	remoteApp := &models.RemoteApp{
		Id: remoteAppId,
	}

	updateCols := make([]string, 0)
	if req.SoftwareId != nil {
		remoteApp.SoftwareId = snowflake.MustParseString(*req.SoftwareId)
		updateCols = append(updateCols, "software_id")
	}

	if req.Name != nil {
		remoteApp.Name = *req.Name
		updateCols = append(updateCols, "name")
	}

	if req.Desc != nil {
		remoteApp.Desc = *req.Desc
		updateCols = append(updateCols, "desc")
	}

	if req.Dir != nil {
		remoteApp.Dir = *req.Dir
		updateCols = append(updateCols, "dir")
	}

	if req.Args != nil {
		remoteApp.Args = *req.Args
		updateCols = append(updateCols, "args")
	}

	if req.Logo != nil {
		remoteApp.Logo = *req.Logo
		updateCols = append(updateCols, "logo")
	}

	if req.DisableGfx != nil {
		remoteApp.DisableGfx = *req.DisableGfx
		updateCols = append(updateCols, "disable_gfx")
	}

	if req.LoginUser != nil {
		remoteApp.LoginUser = *req.LoginUser
		updateCols = append(updateCols, "login_user")
	}

	return remoteApp, updateCols
}

func DeleteRemoteApp(c *gin.Context) {
	logger := trace.GetLogger(c)

	req := new(remoteappmodel.AdminDeleteRequest)
	err := bindDeleteRemoteAppRequest(req, c)
	if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind delete remote app request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminDeleteRemoteAppRequest(req)
	if err = response.BadRequestIfError(c, err, errResp); err != nil {
		logger.Warnf("validate admin delete remote app request failed, %v", err)
		return
	}

	exist, err := dao.DeleteRemoteApp(c, snowflake.MustParseString(*req.RemoteAppId))
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "delete remote app failed")); err != nil {
		logger.Warnf("delete remote app failed, %v", err)
		return
	}
	if !exist {
		err = fmt.Errorf("remoteapp not found")
		_ = response.NotFoundIfError(c, err, response.WrapErrorResp(common.RemoteAppNotFound, err.Error()))
		logger.Warn(err)
		return
	}

	response.RenderJson(nil, c)
}

func bindDeleteRemoteAppRequest(req *remoteappmodel.AdminDeleteRequest, c *gin.Context) error {
	if err := c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	return nil
}
