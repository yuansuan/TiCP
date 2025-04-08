package v20230530

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	copy2 "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/copy"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Copy 复制一个文件/目录(包括子目录和文件)
func (s *Storage) Copy(ctx *gin.Context) {

	userID, accessKey, systemFlag, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "Copy", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &copy2.Request{}
	if err := ctx.BindJSON(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	flag, _, msg := fsutil.ValidateUserIDPath(request.SrcPath)
	if !flag {
		logger.Error(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}

	flag, pathUserID, msg := fsutil.ValidateUserIDPath(request.DestPath)
	if !flag {
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}
	// generate absolute path
	srcAbsPath := filepath.Join(s.rootPath, fsutil.TrimPrefix(request.SrcPath, "/"))

	//check if user has access to the path
	if !systemFlag {
		if !s.CheckPathAccessAndHandleError(accessKey, userID, request.SrcPath, logger, ctx) ||
			!s.CheckPathAccessAndHandleError(accessKey, userID, request.DestPath, logger, ctx) {
			return
		}
	}

	// stat file or directory
	srcFileInfo, err := os.Stat(srcAbsPath)
	if err != nil {
		if os.IsNotExist(err) {
			msg := "file not found"
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusNotFound, commoncode.SrcPathNotFound, msg)
			return
		}
		msg := fmt.Sprintf("stat file error, err: %v", err)
		logger.Errorf(msg)
		common.InternalServerError(ctx, msg)
		return
	}

	//check storage usage
	if !s.Quota.CheckStorageUsageAndHandleError(pathUserID, s.GetTotalSize(logger, srcAbsPath), logger, ctx) {
		return
	}

	// generate absolute path
	destAbsPath := filepath.Join(s.rootPath, fsutil.TrimPrefix(request.DestPath, "/"))

	var destFileInfo os.FileInfo
	if flag, destFileInfo, err = s.HandlePathContainsFileError(ctx, logger, destAbsPath, fmt.Sprintf("dest path contains file, path: %s", request.DestPath)); !flag {
		return
	}

	if err != nil && !os.IsNotExist(err) {
		msg := fmt.Sprintf("stat file or directory error, err: %v,path: %v", err, request.DestPath)
		logger.Errorf(msg)
		common.InternalServerError(ctx, "stat file or directory error")
		return
	}
	if destFileInfo != nil {
		msg := fmt.Sprintf("file or directory alreay exists, path: %v", request.DestPath)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.DestPathExists, msg)
		return
	}

	if srcFileInfo.IsDir() {
		// copy directory recursively
		totalSize, err := s.copyDirectory(srcAbsPath, destAbsPath, logger)
		if err != nil {
			msg := fmt.Sprintf("copy directory error, err: %v", err)
			logger.Errorf(msg)
			common.InternalServerError(ctx, "copy directory error")
			return
		}
		s.OperationLog.InsertOperationLog(logger, ctx, &model.StorageOperationLog{
			UserId:        userID,
			FileName:      filepath.Base(request.DestPath),
			SrcPath:       request.SrcPath,
			DestPath:      request.DestPath,
			FileType:      commoncode.FOLDER,
			OperationType: commoncode.COPY,
			Size:          s.OperationLog.FormatBytes(totalSize),
			CreateTime:    time.Now(),
		})

		common.SuccessResp(ctx, nil)
		return
	}

	// copy single file
	err = s.CopyFile(srcAbsPath, destAbsPath, srcFileInfo.Size(), logger)
	if err != nil {
		msg := fmt.Sprintf("copy file error, srcPath: %s, destPath: %s, err: %v", request.SrcPath, request.DestPath, err)
		logger.Errorf(msg)
		common.InternalServerError(ctx, msg)
	}

	s.OperationLog.InsertOperationLog(logger, ctx, &model.StorageOperationLog{
		UserId:        userID,
		FileName:      filepath.Base(request.DestPath),
		SrcPath:       request.SrcPath,
		DestPath:      request.DestPath,
		FileType:      commoncode.FILE,
		OperationType: commoncode.COPY,
		Size:          s.OperationLog.FormatBytes(srcFileInfo.Size()),
		CreateTime:    time.Now(),
	})

	common.SuccessResp(ctx, nil)
}

func (s *Storage) copyDirectory(srcDir, destDir string, logger *logging.Logger) (int64, error) {
	// create destination directory
	err := os.MkdirAll(destDir, filemode.Directory)
	if err != nil {
		return 0, err
	}

	// read source directory
	fileInfos, err := os.ReadDir(srcDir)
	if err != nil {
		return 0, err
	}

	// copy files and directories
	var totalSize int64
	for _, fileInfo := range fileInfos {
		srcPath := filepath.Join(srcDir, fileInfo.Name())
		destPath := filepath.Join(destDir, fileInfo.Name())

		if fileInfo.IsDir() {
			// if it's a directory, copy recursively
			size, err := s.copyDirectory(srcPath, destPath, logger)
			if err != nil {
				logger.Errorf("copy directory error, srcPath: %s, destPath: %s, err: %v", srcPath, destPath, err)
				return 0, err
			}
			totalSize += size
		} else {
			// if it's a file, copy it
			info, err := fileInfo.Info()
			if err != nil {
				logger.Errorf("get file info failed")
				return 0, err
			}
			err = s.CopyFile(srcPath, destPath, info.Size(), logger)
			if err != nil {
				logger.Errorf("copy file error, srcPath: %s, destPath: %s, err: %v", srcPath, destPath, err)
				return 0, err
			}
			totalSize += info.Size()
		}
	}

	return totalSize, nil
}

func (s *Storage) CopyFile(srcAbsPath, destAbsPath string, fileSize int64, logger *logging.Logger) error {
	// open src
	src, err := os.Open(srcAbsPath)
	defer func(src *os.File) {
		if src == nil {
			return
		}
		err := src.Close()
		if err != nil {
			logger.Errorf("close file error, err: %v", err)
		}
	}(src)
	if err != nil {
		msg := fmt.Sprintf("open file error, err: %v", err)
		logger.Error(msg)
		return err
	}

	// create dest file
	dest, err := fsutil.CreateAndReturnFile(destAbsPath, fileSize, logger)
	defer func(dest *os.File) {
		err := dest.Close()
		if dest == nil {
			return
		}
		if err != nil {
			logger.Errorf("close file error, err: %v", err)
		}
	}(dest)
	if err != nil {
		msg := fmt.Sprintf("create file error, err: %v", err)
		logger.Error(msg)
		return err
	}

	// copy it
	buf := s.pool.Get()
	defer s.pool.Put(buf)

	if _, err := io.CopyBuffer(dest, src, buf.Bytes()); err != nil {
		msg := fmt.Sprintf("copy file error, err: %v", err)
		logger.Error(msg)
		return err
	}

	// sync it
	if err := dest.Sync(); err != nil {
		msg := fmt.Sprintf("sync file error, err: %v", err)
		logger.Error(msg)
		return err
	}

	return nil
}
