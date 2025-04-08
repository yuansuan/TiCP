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
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/hashid"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	sharedDirectoryService "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/service/shareddirectory"
	"go.uber.org/zap"
)

// Create 创建共享目录
func (s *SharedDirectory) Create(ctx *gin.Context) {
	userID, accessKey, _, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "SharedDirectoryCreate", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)
	req := &api.CreateSharedDirectoryRequest{}
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

	ignore := req.IgnoreExisting
	existSharedDirectorys := make([]*model.SharedDirectory, 0, len(req.Paths))
	sharedDirectorys := make([]*model.SharedDirectory, 0, len(req.Paths))
	encodeStrs := make(map[string]string)
	for _, path := range req.Paths {
		flag, pathUser, msg := fsutil.ValidateUserIDPath(path)
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

		// return encode error
		encodeStr, err := hashid.EncodeStr(path)
		if err != nil {
			msg := fmt.Sprintf("encode path failed, err: %v", err)
			logger.Warn(msg)
			common.ErrorResp(ctx, http.StatusInternalServerError, commoncode.InternalServerErrorCode, msg)
			return
		}

		encodeStrs[path] = encodeStr

		mp, err := genSharedDirectoryModel(ctx, logger, path, pathUser)
		if err != nil {
			common.ErrorResp(ctx, http.StatusInternalServerError, commoncode.InternalServerErrorCode, msg)
			return
		}

		exist, sharedDirectory, err := sharedDirectoryService.GetSharedDirectoryInfoByPath(ctx, s.Engine, s.StorageSharedDirectoryDao, path)
		if err != nil {
			msg := fmt.Sprintf("get shared directory info by path failed, err: %v", err)
			logger.Warn(msg)
			common.ErrorResp(ctx, http.StatusInternalServerError, commoncode.InternalServerErrorCode, msg)
			return
		}

		if exist && ignore {
			logger.Infof("shared directory '%s' already exist, will recreate", path)
			existSharedDirectorys = append(existSharedDirectorys, sharedDirectory)
			continue
		}

		if exist {
			msg := fmt.Sprintf("shared directory '%s' already exist", path)
			logger.Warn(msg)
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.SharedDirectoryExisting, msg)
			return
		}

		sharedDirectorys = append(sharedDirectorys, mp)
	}

	shareRegisterAddress := config.GetConfig().ShareRegisterAddress

	for _, sharedDirectory := range sharedDirectorys {
		// dirnfs POST /users/{userName}
		addUserReq := genAddUserReq(sharedDirectory, encodeStrs[sharedDirectory.Path])

		err := s.createSharedDirectory(ctx, addUserReq, sharedDirectory, shareRegisterAddress, logger)
		if err != nil {
			common.ErrorResp(ctx, http.StatusInternalServerError, commoncode.InternalServerErrorCode, err.Error())
			return
		}
	}

	// 重复的共享目录不报错，直接调用创建api并返回已存在的共享目录
	for _, sharedDirectory := range existSharedDirectorys {
		addUserReq := genAddUserReq(sharedDirectory, encodeStrs[sharedDirectory.Path])
		err := s.createSharedDirectoryNoDB(ctx, addUserReq, sharedDirectory, shareRegisterAddress, logger)
		if err != nil {
			common.ErrorResp(ctx, http.StatusInternalServerError, commoncode.InternalServerErrorCode, err.Error())
			return
		}

		sharedDirectorys = append(sharedDirectorys, sharedDirectory)
	}

	common.SuccessResp(ctx, ToResponseSharedDirectorys(sharedDirectorys))
}

func (s *SharedDirectory) createSharedDirectory(ctx *gin.Context, addUserReq dirnfsmodel.AddUserRequest, sharedDirectory *model.SharedDirectory, shareRegisterAddress string, logger *zap.SugaredLogger) error {
	// TODO: dirnfs之后可能添加一个批量添加用户的接口
	err := s.addUser(ctx, addUserReq, sharedDirectory, shareRegisterAddress, logger)
	if err != nil {
		return err
	}

	err = sharedDirectoryService.InsertSharedDirectoryInfo(ctx, s.Engine, s.StorageSharedDirectoryDao, sharedDirectory)
	if err != nil {
		msg := fmt.Sprintf("insert shared directory info failed, err: %v", err)
		logger.Warn(msg)
		return fmt.Errorf(msg)
	}
	return nil
}

func (s *SharedDirectory) createSharedDirectoryNoDB(ctx *gin.Context, addUserReq dirnfsmodel.AddUserRequest, sharedDirectory *model.SharedDirectory, shareRegisterAddress string, logger *zap.SugaredLogger) error {
	err := s.addUser(ctx, addUserReq, sharedDirectory, shareRegisterAddress, logger)
	if err != nil {
		return err
	}

	return nil
}

func (s *SharedDirectory) addUser(ctx *gin.Context, req dirnfsmodel.AddUserRequest, sharedDirectory *model.SharedDirectory, shareRegisterAddress string, logger *zap.SugaredLogger) error {
	resp, err := s.hc.R().
		SetBody(req).
		SetPathParam("username", sharedDirectory.SharedUserName).
		Post(fmt.Sprintf("http://%s/users/{username}", shareRegisterAddress))
	if err != nil {
		msg := fmt.Sprintf("add user failed, err: %v", err)
		logger.Warn(msg)
		return fmt.Errorf(msg)
	}
	if resp.StatusCode() != http.StatusOK {
		msg := fmt.Sprintf("add user failed, status code: %d, body: %s", resp.StatusCode(), resp.String())
		logger.Warn(msg)
		return fmt.Errorf(msg)
	}
	return nil
}

func genAddUserReq(sharedDirectory *model.SharedDirectory, encodeStr string) dirnfsmodel.AddUserRequest {
	return dirnfsmodel.AddUserRequest{
		Username:      sharedDirectory.SharedUserName,
		Password:      sharedDirectory.SharedPassword,
		SubPath:       encodeStr,
		ExcludeUserID: true,
	}
}

func genSharedDirectoryModel(ctx *gin.Context, logger *zap.SugaredLogger, path, pathUser string) (*model.SharedDirectory, error) {
	originalUsername, err := generateUserName()
	if err != nil {
		msg := fmt.Sprintf("generate username failed, err: %v", err)
		logger.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	password, err := generateRandomPassword(DefaultPWLength)
	if err != nil {
		msg := fmt.Sprintf("generate password failed, err: %v", err)
		logger.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	sharedHost := config.GetConfig().SharedHost
	sharedSrc := addPrefixToUserName(originalUsername)

	return &model.SharedDirectory{
		Path:           path,
		UserID:         pathUser,
		SharedUserName: originalUsername, // 数据库存的是原始的username，不包含前缀
		SharedPassword: password,
		SharedHost:     sharedHost,
		SharedSrc:      sharedSrc,
	}, nil
}
