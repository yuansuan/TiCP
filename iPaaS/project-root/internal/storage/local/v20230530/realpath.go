package v20230530

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/realpath"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"net/http"
	"os"
	"path/filepath"
)

// Realpath 根据相对路径转成绝对路径
func (s *Storage) Realpath(ctx *gin.Context) {
	userID, accessKey, _, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "Realpath", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID, "accessKey", accessKey)

	request := &realpath.Request{}
	if err := ctx.BindQuery(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	flag, _, msg := fsutil.ValidateUserIDPath(request.RelativePath)
	if !flag {
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}

	// generate absolute path
	absPath := filepath.Join(s.rootPath, fsutil.TrimPrefix(request.RelativePath, "/"))

	// stat file or directory
	_, err = os.Stat(absPath)
	if err != nil && !os.IsNotExist(err) {
		msg := fmt.Sprintf("stat file or directory error, path: %s, err: %v", absPath, err)
		logger.Errorf(msg)
		common.InternalServerError(ctx, "stat file or directory error")
		return
	}

	data := new(realpath.Data)
	data.RealPath = absPath
	common.SuccessResp(ctx, data)
}
