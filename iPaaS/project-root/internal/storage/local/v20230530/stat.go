package v20230530

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/stat"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"net/http"
	"os"
	"path/filepath"
)

// Stat 获取一个文件或文件夹的信息。
func (s *Storage) Stat(ctx *gin.Context) {

	userID, accessKey, systemFlag, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "Stat", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &stat.Request{}
	if err := ctx.BindQuery(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	flag, _, msg := fsutil.ValidateUserIDPath(request.Path)
	if !flag {
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}

	//check if user has access to the path
	if !systemFlag && !s.CheckPathAccessAndHandleError(accessKey, userID, request.Path, logger, ctx) {
		return
	}

	// generate absolute path
	absPath := filepath.Join(s.rootPath, fsutil.TrimPrefix(request.Path, "/"))

	// stat file or directory
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			msg := "file or directory not found,path: " + request.Path
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusNotFound, commoncode.PathNotFound, msg)
			return
		}
		msg := fmt.Sprintf("stat file error, err: %v", err)
		logger.Errorf(msg)
		common.InternalServerError(ctx, "stat file error")
		return
	}

	data := stat.Data{
		File: v20230530.ToRespFileInfo(fileInfo),
	}
	common.SuccessResp(ctx, data)
}
