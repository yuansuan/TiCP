package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/shared_directory/api"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	sharedDirectoryService "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/service/shareddirectory"
)

// List 列出共享目录
func (s *SharedDirectory) List(ctx *gin.Context) {
	userID, accessKey, _, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "SharedDirectoryList", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)
	req := &api.ListSharedDirectoryRequest{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Warn(msg)
		common.InvalidParams(ctx, err.Error())
		return
	}

	pathPrefix := req.PathPrefix

	flag, _, msg := fsutil.ValidateUserIDPath(pathPrefix)
	if !flag {
		logger.Warn(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}

	// Check if user has access to the path
	// No system level APIs are currently required
	// if !systemFlag && !s.CheckPathAccessAndHandleError(accessKey, userID, request.Path, logger, ctx) {
	accessFlag := s.CheckPathAccessAndHandleError(accessKey, userID, pathPrefix, logger, ctx)
	if !accessFlag {
		return
	}

	// 根据pathUser查询
	sharedDirectorys, err := sharedDirectoryService.ListSharedDirectoryInfoByPathPrefix(ctx, s.Engine, s.StorageSharedDirectoryDao, pathPrefix)
	if err != nil {
		msg := fmt.Sprintf("list shared directory info by path failed, err: %v", err)
		logger.Warn(msg)
		common.ErrorResp(ctx, http.StatusInternalServerError, commoncode.InternalServerErrorCode, msg)
		return
	}

	resp := ToResponseSharedDirectorys(sharedDirectorys)

	common.SuccessResp(ctx, resp)
}
