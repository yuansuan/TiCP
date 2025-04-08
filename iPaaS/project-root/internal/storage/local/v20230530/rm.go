package v20230530

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/rm"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	maxRetries    = 10
	retryInterval = time.Millisecond * 100
)

// Rm 删除一个文件或文件夹。
func (s *Storage) Rm(ctx *gin.Context) {

	userID, accessKey, systemFlag, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "Rm", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &rm.Request{}
	if err := ctx.BindJSON(request); err != nil {
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

	// stat file
	var fileInfo os.FileInfo
	if !request.IgnoreNotExist {
		fileInfo, err = os.Lstat(absPath)
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
	}
	fileType := commoncode.FILE
	var totalSize int64
	if fileInfo != nil && fileInfo.IsDir() {
		fileType = commoncode.FOLDER
		err := filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("Error accessing path %s: %s\n", path, err)
				return nil
			}
			if info.IsDir() {
				return nil
			}
			totalSize += info.Size()
			return nil
		})

		if err != nil {
			fmt.Printf("Error walking directory: %s\n", err)
			return
		}
	} else if fileInfo != nil {
		totalSize = fileInfo.Size()
	}

	if err = removeFileOrDirectoryWithRetry(absPath, request.Path, ctx, logger); err != nil {
		return
	}

	s.OperationLog.InsertOperationLog(logger, ctx, &model.StorageOperationLog{
		UserId:        userID,
		FileName:      filepath.Base(request.Path),
		SrcPath:       request.Path,
		FileType:      fileType,
		OperationType: commoncode.DELETE,
		Size:          s.OperationLog.FormatBytes(totalSize),
		CreateTime:    time.Now(),
	})

	logger.Infof("remove file or directory success, path: %s", request.Path)

	common.SuccessResp(ctx, nil)
}

func removeFileOrDirectoryWithRetry(absPath string, requestPath string, ctx *gin.Context, logger *logging.Logger) error {
	msg := fmt.Sprintf("remove file or directory error, path: %s", requestPath)
	var err error
	for i := 0; i <= maxRetries; i++ {
		err = os.RemoveAll(absPath)
		if err == nil {
			return nil
		}

		logger.Infof(msg+",retry "+strconv.Itoa(i)+"time,error:", err.Error())

		if i != maxRetries {
			time.Sleep(retryInterval)
		}
	}

	logger.Errorf(msg+",error:", err.Error())
	common.InternalServerError(ctx, "remove file or directory error")
	return errors.New("remove file or directory error")
}
