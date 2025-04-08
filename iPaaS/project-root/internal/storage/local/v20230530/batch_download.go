package v20230530

import (
	"archive/zip"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/batchDownload"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
)

func (s *Storage) BatchDownload(ctx *gin.Context) {

	userID, accessKey, systemFlag, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "BatchDownload", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &batchDownload.Request{}
	if err := ctx.BindJSON(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, err.Error())
		return
	}

	if len(request.Paths) == 0 {
		msg := "paths is empty"
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}

	if len(strings.TrimSpace(request.FileName)) == 0 {
		msg := "file name is empty"
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidFileName, msg)
		return
	}

	for _, path := range request.Paths {
		flag, _, msg := fsutil.ValidateUserIDPath(path)
		if !flag {
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
			return
		}
		//check if user has access to the path
		if !systemFlag && !s.CheckPathAccessAndHandleError(accessKey, userID, path, logger, ctx) {
			return
		}
	}

	pathsMap := make(map[string]struct{})
	emptyDirs := make([]string, 0)
	needZip := true
	// recursively get all files
	for _, path := range request.Paths {
		// generate absolute path
		absPath := filepath.Join(s.rootPath, fsutil.TrimPrefix(path, "/"))
		// recursively get all files if the path is a directory
		p, fileInfo, err := fsutil.ReadFinalPath(absPath)
		absPath = p
		if err == nil {
			if fileInfo.IsDir() {
				if err = fsutil.GetAllFileAndEmptyDir(absPath, pathsMap, &emptyDirs, logger, ctx, true); err != nil {
					return
				}
			} else {
				pathsMap[absPath] = struct{}{}
				// if only one file and request.path is file, return the file directly
				if len(request.Paths) == 1 {
					needZip = false
				}
			}
		} else {
			if os.IsNotExist(err) {
				msg := fmt.Sprintf("path not found, path: %s", path)
				logger.Info(msg)
				common.ErrorResp(ctx, http.StatusNotFound, commoncode.PathNotFound, msg)
				return
			}
			msg := fmt.Sprintf("Lstat error, path: %s, err: %v", path, err)
			logger.Error(msg)
			common.InternalServerError(ctx, "lstat error")
			return
		}
	}

	var paths []string
	for path := range pathsMap {
		paths = append(paths, path)
	}

	//if only one file and isCompress is false, return the file directly
	if !request.IsCompress && len(paths) == 1 && len(emptyDirs) == 0 {
		needZip = false
	}

	// check if target file is a zip file
	if needZip && filepath.Ext(request.FileName) != ".zip" {
		msg := fmt.Sprintf("unsupported compress file type, fileName: %s", request.FileName)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.UnsupportedCompressFileType, msg)
		return
	}

	ctx.Writer.Header().Set("Content-Type", "application/zip")
	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%v"`, request.FileName))

	basePath := filepath.Join(s.rootPath, userID)
	// if only one file, return the file directly
	if !needZip {
		relPath, err := filepath.Rel(basePath, paths[0])
		if err != nil {
			msg := fmt.Sprintf("get relative path error, err: %v", err)
			logger.Error(msg)
			common.InternalServerError(ctx, "get relative path error")
			return
		}
		fileInfo, err := os.Lstat(paths[0])
		if err != nil {
			msg := fmt.Sprintf("lstat error, path: %s, err: %v", relPath, err)
			logger.Errorf(msg)
			common.InternalServerError(ctx, "lstat error")
			return
		}
		file, err := os.Open(paths[0])
		if err != nil {
			msg := fmt.Sprintf("open file error, path: %s, err: %v", relPath, err)
			logger.Errorf(msg)
			common.InternalServerError(ctx, "open file error")
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				logger.Errorf("close file error, err: %v", err)
			}
		}(file)

		fileExt := filepath.Ext(paths[0])
		contentType := mime.TypeByExtension(fileExt)
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		ctx.Writer.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
		ctx.Writer.Header().Set("Content-Type", contentType)

		s.OperationLog.InsertOperationLog(logger, ctx, &model.StorageOperationLog{
			UserId:        userID,
			FileName:      filepath.Base(paths[0]),
			SrcPath:       strings.TrimPrefix(paths[0], s.rootPath),
			FileType:      commoncode.FILE,
			OperationType: commoncode.DOWNLOAD,
			Size:          s.OperationLog.FormatBytes(fileInfo.Size()),
			CreateTime:    time.Now(),
		})

		http.ServeContent(ctx.Writer, ctx.Request, request.FileName, fileInfo.ModTime(), file)
		return
	}

	if request.BasePath != "" {
		flag, _, msg := fsutil.ValidateUserIDPath(request.BasePath)
		if !flag {
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
			return
		}
		basePath = filepath.Join(s.rootPath, fsutil.TrimPrefix(request.BasePath, "/"))
	}

	zipWriter := zip.NewWriter(ctx.Writer)
	defer func() {
		err := zipWriter.Close()
		if err != nil {
			logger.Errorf("close zipWriter error, err: %v", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	//add file to zip
	go func() {
		defer wg.Done()
		allPaths := append(paths, emptyDirs...)
		err := s.createZipFiles(basePath, zipWriter, logger, allPaths...)
		if err != nil {
			msg := fmt.Sprintf("add file dir to zip file error, err: %v", err)
			logger.Errorf(msg)
			common.InternalServerError(ctx, "create zip file error")
			return
		}
	}()

	s.OperationLog.InsertOperationLog(logger, ctx, &model.StorageOperationLog{
		UserId:        userID,
		FileName:      fmt.Sprintf(request.FileName),
		SrcPath:       strings.Join(request.Paths, ","),
		FileType:      commoncode.Batch,
		OperationType: commoncode.DOWNLOAD,
		Size:          "",
		CreateTime:    time.Now(),
	})

	wg.Wait()

}
