package v20230530

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/ls"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

func (s *Storage) Ls(ctx *gin.Context) {

	userID, accessKey, systemFlag, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "Ls", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &ls.Request{}
	if err := ctx.BindQuery(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, err.Error())
		return
	}

	//FilterRegexpList参数放在请求body里
	requestBody := &ls.Request{}
	if err := ctx.ShouldBindJSON(requestBody); err != nil && err.Error() != "EOF" {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, err.Error())
		return
	}

	if request.PageOffset < 0 {
		msg := fmt.Sprintf("invalid page offset, page offset should be greater than or equal to 0, got: %v", request.PageOffset)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPageOffset, msg)
		return
	}

	if request.PageSize < 1 || request.PageSize > 1000 {
		msg := fmt.Sprintf("invalid page size, page size should be in [1, 1000], got: %v", request.PageSize)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPageSize, msg)
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
	filePath := filepath.Join(s.rootPath, fsutil.TrimPrefix(request.Path, "/"))

	//check if path exists and is a directory
	fileInfo, err := os.Lstat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			//if path is user root path, create user folder and return empty list
			if request.Path == "/"+userID+"/" {
				if err := os.MkdirAll(filePath, filemode.Directory); err != nil {
					msg := fmt.Sprintf("mkdir error, path: %s, err: %v", request.Path, err)
					logger.Error(msg)
					common.InternalServerError(ctx, "internal server error")
					return
				}
				common.SuccessResp(ctx, &ls.Data{
					Files:      []*v20230530.FileInfo{},
					Total:      0,
					NextMarker: -1,
				})
				return
			} else {
				msg := fmt.Sprintf("directory path not found, path: %s", request.Path)
				logger.Info(msg)
				common.ErrorResp(ctx, http.StatusNotFound, commoncode.PathNotFound, msg)
				return
			}
		}
		msg := fmt.Sprintf("Lstat error, path: %s, err: %v", request.Path, err)
		logger.Error(msg)
		common.InternalServerError(ctx, "Lstat error")
		return
	}
	if !fileInfo.IsDir() {
		if fileInfo.Mode()&os.ModeSymlink != 0 {
			logger.Debugf("path %s is a symlink", filePath)
			sourceFilePath, err := os.Readlink(filePath)
			logger.Debugf("source file path: %s", sourceFilePath)
			if err != nil {
				msg := fmt.Sprintf("read symlink error, err: %v", err)
				logger.Errorf(msg)
				common.InternalServerError(ctx, "internal server error")
				return
			}
			if filepath.IsAbs(sourceFilePath) {
				filePath = sourceFilePath
			} else {
				filePath = filepath.Join(filepath.Dir(filePath), sourceFilePath)
			}
		} else {
			msg := fmt.Sprintf("path should be a directory, got: %s", request.Path)
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
			return
		}
	}
	filterRegexpList := make([]*regexp.Regexp, 0)

	// 优先处理 requestBody.FilterRegexpList
	if requestBody.FilterRegexpList != nil && len(requestBody.FilterRegexpList) > 0 {
		for _, filterRegexp := range requestBody.FilterRegexpList {
			regexpObj, err := regexp.Compile(filterRegexp)
			if err != nil {
				msg := fmt.Sprintf("invalid filter regexp, err: %v, filterRegexp: %v", err, filterRegexp)
				logger.Info(msg)
				common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidRegexp, "invalid filter regexp")
				return
			}
			filterRegexpList = append(filterRegexpList, regexpObj)
		}
	}

	// 再处理单一的 request.FilterRegexp
	if request.FilterRegexp != "" {
		regexpObj, err := regexp.Compile(request.FilterRegexp)
		if err != nil {
			msg := fmt.Sprintf("invalid filter regexp, err: %v, filterRegexp: %v", err, request.FilterRegexp)
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidRegexp, "invalid filter regexp")
			return
		}
		filterRegexpList = append(filterRegexpList, regexpObj)
	}

	// 确保 filterRegexpList 非空
	if len(filterRegexpList) == 0 {
		logger.Info("no filter regexp provided, defaulting to all files")
	}

	// 调用 fsutil.LsPage 时正确传递 filterRegexpList
	infos, nextMarker, total, err := fsutil.LsPage(filePath, filterRegexpList, request.PageOffset, request.PageSize, logger)
	if err != nil {
		msg := fmt.Sprintf("lsPage failed, err: %v", err)
		logger.Error(msg)
		common.InternalServerError(ctx, "lsPage failed")
		return
	}

	// 处理结果并返回
	resp := new(ls.Data)
	resp.Files = v20230530.ToRespFileInfos(infos)
	resp.Total = total
	resp.NextMarker = nextMarker

	common.SuccessResp(ctx, resp)

}
