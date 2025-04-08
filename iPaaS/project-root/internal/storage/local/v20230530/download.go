package v20230530

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/download"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// RangeHeaderPattern range header pattern
	RangeHeaderPattern = `bytes=(\d+)-(\d+)`
)

// Download 用于下载一个文件。
// note: 该接口在csp不同的场景下，理解的角度会有不同。
// 在cloud storage中，download 接口就是理解上的下载，从 cloud-storage 下载文件到客户端。
// 但是在hpc-storage中，csp 作业执行中，download 接口被用于从 cloud-storage 下载文件到 hpc-storage。
// 因此，在hpc-storage 场景中，csp 使用 download 下载文件是不会也不应该做quota 限制，所以会出现csp 账号quota 超出配置额度上限的情况。
func (s *Storage) Download(ctx *gin.Context) {

	userID, accessKey, systemFlag, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "Download", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &download.Request{}
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
			msg := fmt.Sprintf("file or directory not found, path: %v", request.Path)
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusNotFound, commoncode.PathNotFound, msg)
			return
		}
		msg := fmt.Sprintf("stat file error, err: %v , path: %v", err, request.Path)
		logger.Errorf(msg)
		common.InternalServerError(ctx, msg)
		return
	}
	if fileInfo.IsDir() {
		msg := fmt.Sprintf("path should be a file, got: %s", request.Path)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}

	//deal with symlink
	if flag, absPath, fileInfo = s.HandleSymbolicLink(absPath, logger, fileInfo, ctx); !flag {
		return
	}

	start := int64(0)
	length := fileInfo.Size()
	if request.Range != "" {
		startOffset, endOffset, err := parseByteRange(request.Range, fileInfo.Size())
		if err != nil {
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidRange, "parse range error")
			return
		}
		start = startOffset
		length = endOffset - startOffset
		ctx.Writer.Header().Set("Accept-Ranges", "bytes")
		ctx.Writer.Header().Set("Range", fmt.Sprintf("bytes=%d-%d", start, start+length-1))
		ctx.Writer.Header().Set("Content-Length", strconv.FormatInt(length, 10))
		ctx.Writer.WriteHeader(http.StatusPartialContent)
	} else {
		ctx.Writer.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
		ctx.Writer.WriteHeader(http.StatusOK)
	}

	fileExt := filepath.Ext(absPath)
	contentType := mime.TypeByExtension(fileExt)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	filename := filepath.Base(absPath)
	ctx.Writer.Header().Set("Content-Type", contentType)
	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%v"`, filename))

	file, err := os.Open(absPath)
	if err != nil {
		msg := fmt.Sprintf("open file error, err: %v , path: %v", err, strings.TrimPrefix(absPath, s.rootPath+"/"))
		logger.Errorf(msg)
		common.InternalServerError(ctx, "open file error")
		return
	}
	defer func() { _ = file.Close() }()
	httpFile := io.NewSectionReader(file, start, length)

	s.OperationLog.InsertOperationLog(logger, ctx, &model.StorageOperationLog{
		UserId:        userID,
		FileName:      filepath.Base(request.Path),
		SrcPath:       request.Path,
		FileType:      commoncode.FILE,
		OperationType: commoncode.DOWNLOAD,
		Size:          s.OperationLog.FormatBytes(length),
		CreateTime:    time.Now(),
	})

	http.ServeContent(ctx.Writer, ctx.Request, filename, fileInfo.ModTime(), httpFile)
}

func parseByteRange(byteRange string, fileSize int64) (int64, int64, error) {
	pattern := RangeHeaderPattern

	r, err := regexp.Compile(pattern)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to compile regexp %s: %v", pattern, err)
	}

	matches := r.FindStringSubmatch(byteRange)
	if len(matches) != 3 {
		return 0, 0, fmt.Errorf("invalid byte range: %s", byteRange)
	}

	start, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start offset: %v", err)
	}
	if start < 0 {
		return 0, 0, fmt.Errorf("start offset can not be negative")
	}

	end, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid end offset: %v", err)
	}
	if end > fileSize {
		return 0, 0, fmt.Errorf("end offset can not be greater than file size")
	}
	// 超算上有0KB的文件，允许下载空文件
	if end < start || (fileSize != 0 && end == start) {
		return 0, 0, fmt.Errorf("end offset must be greater than start offset")
	}

	return start, end, nil
}
