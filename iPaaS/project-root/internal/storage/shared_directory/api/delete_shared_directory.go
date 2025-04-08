package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/shared_directory/api"
	dirnfsmodel "github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_dirnfs/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	sharedDirectoryService "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/service/shareddirectory"
	"go.uber.org/zap"
)

// Delete 删除共享目录
func (s *SharedDirectory) Delete(ctx *gin.Context) {
	userID, accessKey, _, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "SharedDirectoryDelete", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)
	req := &api.DeleteSharedDirectoryRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Warn(msg)
		common.InvalidParams(ctx, err.Error())
		return
	}

	if len(req.Paths) == 0 {
		msg := "paths can not be empty"
		logger.Warn(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}

	ignore := req.IgnoreNonexistent
	sharedDirectorys := make([]*model.SharedDirectory, 0, len(req.Paths))
	for _, path := range req.Paths {
		flag, _, msg := fsutil.ValidateUserIDPath(path)
		if !flag {
			logger.Warn(msg)
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
			return
		}

		// Check if user has access to the path
		// No system level APIs are currently required
		// if !systemFlag && !s.CheckPathAccessAndHandleError(accessKey, userID, request.Path, logger, ctx) {
		accessFlag := s.CheckPathAccessAndHandleError(accessKey, userID, path, logger, ctx)
		if !accessFlag {
			return
		}

		// 根据path查询数据库
		exist, sharedDirectory, err := sharedDirectoryService.GetSharedDirectoryInfoByPath(ctx, s.Engine, s.StorageSharedDirectoryDao, path)
		if err != nil {
			msg := fmt.Sprintf("get shared directory info by path failed, err: %v", err)
			logger.Warn(msg)
			common.ErrorResp(ctx, http.StatusInternalServerError, commoncode.InternalServerErrorCode, msg)
			return
		}

		if !exist && ignore {
			logger.Infof("shared directory '%s' not exist, ignore", path)
			continue
		}

		if !exist {
			msg := fmt.Sprintf("shared directory '%s' not exist", path)
			logger.Warn(msg)
			common.ErrorResp(ctx, http.StatusNotFound, commoncode.SharedDirectoryNonexistent, msg)
			return
		}

		sharedDirectorys = append(sharedDirectorys, sharedDirectory)
	}

	for _, sharedDirectory := range sharedDirectorys {
		// dirnfs DELETE /users/{userName}
		deleteUserReq := genDeleteUserReq(sharedDirectory)

		shareRegisterAddress := config.GetConfig().ShareRegisterAddress
		err = s.deleteSharedDirectory(ctx, *deleteUserReq, sharedDirectory, shareRegisterAddress, logger)
		if err != nil {
			common.ErrorResp(ctx, http.StatusInternalServerError, commoncode.InternalServerErrorCode, err.Error())
			return
		}

	}

	common.SuccessResp(ctx, nil)
}

func (s *SharedDirectory) deleteSharedDirectory(ctx *gin.Context, deleteUserReq dirnfsmodel.DeleteUserRequest, sharedDirectory *model.SharedDirectory, shareRegisterAddress string, logger *zap.SugaredLogger) error {
	// TODO: dirnfs之后可能添加一个批量删除用户的接口
	err := s.deleteUser(ctx, deleteUserReq, sharedDirectory, shareRegisterAddress, logger)
	if err != nil {
		return err
	}

	err = sharedDirectoryService.DeleteSharedDirectoryInfo(ctx, s.Engine, s.StorageSharedDirectoryDao, sharedDirectory.Path)
	if err != nil {
		msg := fmt.Sprintf("delete shared directory info failed, err: %v", err)
		logger.Warn(msg)
		return fmt.Errorf(msg)
	}

	return nil
}

func (s *SharedDirectory) deleteUser(ctx *gin.Context, req dirnfsmodel.DeleteUserRequest, sharedDirectory *model.SharedDirectory, shareRegisterAddress string, logger *zap.SugaredLogger) error {
	resp, err := s.hc.R().
		SetBody(req).
		SetPathParam("username", sharedDirectory.SharedUserName).
		Delete(fmt.Sprintf("http://%s/users/{username}", shareRegisterAddress))
	if err != nil {
		msg := fmt.Sprintf("add user failed, err: %v", err)
		logger.Warn(msg)
		return fmt.Errorf(msg)
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNotFound {
		msg := fmt.Sprintf("add user failed, status code: %d, body: %s", resp.StatusCode(), resp.String())
		logger.Warn(msg)
		return fmt.Errorf(msg)
	}

	return nil
}

func genDeleteUserReq(sharedDirectory *model.SharedDirectory) *dirnfsmodel.DeleteUserRequest {
	deleteUserReq := &dirnfsmodel.DeleteUserRequest{
		Username: sharedDirectory.SharedUserName,
	}
	return deleteUserReq
}
