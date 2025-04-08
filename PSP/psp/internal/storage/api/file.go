package api

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/user"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/collectionutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"
)

// PreUpload
//
//	@Summary		文件预上传接口
//	@Description	文件预上传接口
//	@Tags			存储-文件
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			path		formData	string					true	"路径"		default("/")
//	@Param			file_size	formData	int						true	"文件大小"		default(0)
//	@Param			cross		formData	bool					false	"是否跨越用户目录"	default(false)
//	@Param			is_cloud	formData	bool					false	"是否云端"		default(false)
//	@Param			user_name	formData	string					false	"用户名"		default("")
//	@Success		200			{object}	dto.PreUploadResponse	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/storage/preUpload [post]
func (s *RouteService) PreUpload(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req dto.PreUploadRequest
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	fileService, err := s.GetFileService()
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRead)
		return
	}

	if strutil.IsEmpty(req.UserName) {
		req.UserName = ginutil.GetUserName(ctx)
	}

	preUploadResp, err := fileService.PreUpload(ctx, &req)
	if err != nil {
		logger.Errorf("pre upload file err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileUploadInit)
		return
	}

	ginutil.Success(ctx, preUploadResp)
}

// Upload
//
//	@Summary		文件分片上传接口
//	@Description	文件分片上传接口
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.UploadRequest	true	"入参"
//	@Success		200
//	@Router			/storage/upload [post]
func (s *RouteService) Upload(ctx *gin.Context) {

	logger := logging.GetLogger(ctx)

	var req dto.UploadRequest
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}
	slice := ctx.Request.MultipartForm.File["slice"]
	if len(slice) == 0 {
		ginutil.Error(ctx, errcode.ErrFileFailToGetSlice, "failed to get file from request, file slice is empty")
		return
	}

	fileService, err := s.GetFileService()
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRead)
		return
	}

	if strutil.IsEmpty(req.UserName) {
		req.UserName = ginutil.GetUserName(ctx)
	}

	err = fileService.Upload(ctx, &req, slice)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileUploadSlice)
		return
	}

	if req.Finish {
		_ = oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_FILE_MANAGER, fmt.Sprintf("用户%v上传文件[%v]", req.UserName, req.Path))
	}

	ginutil.Success(ctx, nil)
}

// BatchDownloadPre
//
//	@Summary		批量文件预下载
//	@Description	批量文件预下载接口
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.BatchDownloadPreRequest	true	"入参"
//	@Success		200		{string}	string
//	@Router			/storage/batchDownloadPre [post]
func (s *RouteService) BatchDownloadPre(ctx *gin.Context) {
	var req dto.BatchDownloadPreRequest
	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	if len(req.FilePaths) == 0 || req.FileName == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	fileService, err := s.GetFileService()
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRead)
		return
	}

	if strutil.IsEmpty(req.UserName) {
		req.UserName = ginutil.GetUserName(ctx)
	}

	token, err := fileService.BatchDownloadPre(ctx, &req)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileUploadSlice)
		return
	}

	ginutil.Success(ctx, token)
}

// BatchDownload
//
//	@Summary		批量下载文件
//	@Description	批量下载文件接口
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.BatchDownloadRequest	true	"请求参数"
//	@Success		200
//	@Router			/storage/batchDownload [get]
func (s *RouteService) BatchDownload(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.BatchDownloadRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	fileService, err := s.GetFileService()
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileBatchDownload)
		return
	}

	userName := ginutil.GetUserName(ctx)
	err = fileService.BatchDownload(ctx, userName, req.Token, req.IsCloud)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileBatchDownload)
		return
	}
}

// Compress
//
//	@Summary		压缩文件
//	@Description	压缩文件接口
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.CompressRequest	true	"请求参数"
//	@Success		200
//	@Router			/storage/compress [post]
func (s *RouteService) Compress(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.CompressRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	fileService, err := s.GetFileService()
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRead)
		return
	}

	if len(req.SrcPaths) == 0 || req.DstPath == "" {
		logger.Error("Parameter error")
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrInvalidParam)
		return
	}

	submit, err := fileService.CompressSubmit(ctx)
	if !submit {
		logger.Error("Executed during task and cannot be submitted")
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileSubmitCompressTask)
		return
	}

	data, err := fileService.Compress(ctx, req)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileCompress)
		return
	}

	ginutil.Success(ctx, data)
}

// CompressTasks
//
//	@Summary		压缩任务信息
//	@Description	压缩任务信息接口
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.CompressTasksRequest	true	"请求参数"
//	@Success		200
//	@Router			/storage/compressTasks [get]
func (s *RouteService) CompressTasks(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.CompressTasksRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	fileService, err := s.GetFileService()
	if err != nil {
		logger.Errorf("get fileservice failed: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRead)
		return
	}

	task, err := fileService.CompressTasks(ctx)
	if err != nil {
		logger.Errorf("get task failed: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileGetCompressTasks)
		return
	}

	ginutil.Success(ctx, task)
}

// Copy
//
//	@Summary		复制接口
//	@Description	复制接口
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.CopyRequest	true	"入参"
//	@Success		200		{string}	key				"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/storage/copy [put]
func (s *RouteService) Copy(ctx *gin.Context) {
	var req dto.CopyRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	if strutil.IsEmpty(req.UserName) {
		req.UserName = ginutil.GetUserName(ctx)
	}

	fileService, err := s.GetFileService()
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailCopy)
		return
	}

	key, err := fileService.Copy(ctx, &req)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailCopy)
		return
	}

	allSrcPaths := collectionutil.MergeSlice(req.SrcFilePaths, req.SrcDirPaths)
	content := fmt.Sprintf("用户%v复制文件%v => [%v]", req.UserName, allSrcPaths, req.DstPath)
	if req.Cross && !strings.HasPrefix(allSrcPaths[0], req.UserName) {
		pathSlice := strings.Split(allSrcPaths[0], "/")
		filepath.Join("/", ".", req.DstPath)
		content = fmt.Sprintf("用户%v保存%v发送的文件[%v] => [%v]", req.UserName, pathSlice[0], pathSlice[len(pathSlice)-1], strings.Replace(req.DstPath, req.UserName, "home", 1))
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_FILE_MANAGER, content)
	ginutil.Success(ctx, key)
}

// CopyStatus
//
//	@Summary		复制接口
//	@Description	复制接口
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			copyKey	query		bool	false	"key"	default(false)
//	@Success		200		{string}	key		"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/storage/copyStatus [get]
func (s *RouteService) CopyStatus(ctx *gin.Context) {
	key := ctx.Query("key")

	if strutil.IsEmpty(key) {
		http.Errf(ctx, errcode.ErrInvalidParam, "key can't empty")
		return
	}

	copyStatus, err := s.LocalFileService.GetCopyStatus(ctx, key)

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailCopy)
		return
	}

	ginutil.Success(ctx, copyStatus)
}

// List
//
//	@Summary		文件列表接口
//	@Description	文件列表接口
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.ListRequest		true	"入参"
//	@Success		200		{array}	dto.ListResponse	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/storage/list [post]
func (s *RouteService) List(ctx *gin.Context) {
	var req dto.ListRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	if strutil.IsEmpty(req.UserName) {
		req.UserName = ginutil.GetUserName(ctx)
	}

	fileService, err := s.GetFileService()
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileList)
		return
	}

	responses := make([]*dto.ListResponse, 0)

	// Check the user home path existence
	err = fileService.CheckUserHomePath(ctx, req.UserName, req.Cross)
	if err != nil {
		ginutil.Success(ctx, responses)
		return
	}

	err = s.LocalFileService.CheckSharePath(ctx, req.UserName, req.Cross)
	if err != nil {
		ginutil.Success(ctx, responses)
		return
	}

	files, err := fileService.List(ctx, req.UserName, req.Path, req.Page, req.Cross, req.ShowHideFile, req.FilterRegexpList)
	if err != nil {
		ginutil.Success(ctx, responses)
		return
	}

	files, err = fileService.CheckOnlyReadDir(ctx, req.Path, files)
	if err != nil {
		ginutil.Success(ctx, responses)
		return
	}

	ginutil.Success(ctx, util.ConvertFileRsps(files))
}

// ListOfRecur
//
//	@Summary		文件列表接口(递归)
//	@Description	文件列表接口(递归)
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.ListRecurRequest	true	"入参"
//	@Success		200		{array}	dto.ListResponse		"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/storage/listOfRecur [post]
func (s *RouteService) ListOfRecur(ctx *gin.Context) {
	var req dto.ListRecurRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	if strutil.IsEmpty(req.UserName) {
		req.UserName = ginutil.GetUserName(ctx)
	}

	fileService, err := s.GetFileService()
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRead)
		return
	}

	responses := make([]*dto.ListResponse, 0)

	// Check the user home path existence
	if !fileService.Exist(ctx, req.UserName, false, ".") {
		err := fileService.CreateDir(ctx, req.UserName, ".", req.Cross)
		if err != nil {
			ginutil.Success(ctx, responses)
			return
		}
	}

	files, err := fileService.ListOfRecur(ctx, req.UserName, req.Paths, req.Cross, req.ShowHideFile, req.FilterRegexpList)

	if err != nil {
		ginutil.Success(ctx, responses)
		return
	}

	ginutil.Success(ctx, util.ConvertFileRsps(files))
}

// Get
//
//	@Summary		文件查询接口
//	@Description	文件查询接口
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.GetRequest		true	"入参"
//	@Success		200		{array}	dto.ListResponse	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/storage/get [post]
func (s *RouteService) Get(ctx *gin.Context) {
	var req dto.GetRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	fileService, err := s.GetFileService()
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRead)
		return
	}

	if strutil.IsEmpty(req.UserName) {
		req.UserName = ginutil.GetUserName(ctx)
	}

	var files []*dto.File
	for _, path := range req.Paths {
		fileTmp, err := fileService.Get(ctx, req.UserName, path, req.Cross)
		if err != nil {
			errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailGet)
			return
		}
		files = append(files, fileTmp)
	}

	ginutil.Success(ctx, util.ConvertFileRsps(files))
}

// Rename
//
//	@Summary		文件重命名接口
//	@Description	文件重命名接口
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.RenameRequest	true	"入参"
//	@Success		200
//	@Router			/storage/rename [put]
func (s *RouteService) Rename(ctx *gin.Context) {
	var req dto.RenameRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	fileService, err := s.GetFileService()
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRead)
		return
	}

	if strutil.IsEmpty(req.UserName) {
		req.UserName = ginutil.GetUserName(ctx)
	}

	err = fileService.Rename(ctx, req.UserName, req.Path, req.NewPath, req.Overwrite, req.Cross)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRename)
		return
	}
	_ = oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_FILE_MANAGER, fmt.Sprintf("用户%v重命名文件[%v] => [%v]", req.UserName, req.Path, req.NewPath))

	ginutil.Success(ctx, nil)
}

// Move
//
//	@Summary		文件移动接口
//	@Description	文件移动接口
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.MoveRequest	true	"入参"
//	@Success		200
//	@Router			/storage/move [put]
func (s *RouteService) Move(ctx *gin.Context) {
	var req dto.MoveRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	fileService, err := s.GetFileService()
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRead)
		return
	}

	if strutil.IsEmpty(req.UserName) {
		req.UserName = ginutil.GetUserName(ctx)
	}

	err = fileService.Move(ctx, req.UserName, req.Cross, req.Overwrite, filepath.Join("/", req.DstPath), req.SrcPaths...)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailMove)
		return
	}
	_ = oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_FILE_MANAGER, fmt.Sprintf("用户%v移动文件%v至[%v]", req.UserName, req.SrcPaths, req.DstPath))

	ginutil.Success(ctx, nil)

}

// CreateDir
//
//	@Summary		文件夹创建接口
//	@Description	文件夹创建接口
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.CreateDirRequest	true	"入参"
//	@Success		200
//	@Router			/storage/createDir [post]
func (s *RouteService) CreateDir(ctx *gin.Context) {
	var req dto.CreateDirRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	fileService, err := s.GetFileService()
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRead)
		return
	}

	if strutil.IsEmpty(req.UserName) {
		req.UserName = ginutil.GetUserName(ctx)
	}

	err = fileService.CreateDir(ctx, req.UserName, req.Path, req.Cross)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailMkdir)
		return
	}
	_ = oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_FILE_MANAGER, fmt.Sprintf("用户%v新建文件夹[%v]", req.UserName, req.Path))

	ginutil.Success(ctx, nil)

}

// Remove
//
//	@Summary		文件删除接口
//	@Description	文件删除接口
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.RemoveRequest	true	"入参"
//	@Success		200
//	@Router			/storage/remove [delete]
func (s *RouteService) Remove(ctx *gin.Context) {
	var req dto.RemoveRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	fileService, err := s.GetFileService()
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRead)
		return
	}

	if strutil.IsEmpty(req.UserName) {
		req.UserName = ginutil.GetUserName(ctx)
	}

	for _, p := range req.Paths {
		path := filepath.Clean(p)

		// Check the path existence
		if !fileService.Exist(ctx, req.UserName, req.Cross, path) {
			errcode.ResolveErrCodeMessage(ctx, nil, errcode.ErrFileNotExist)
			return
		}

		// Remove the file
		err := fileService.Remove(ctx, req.UserName, path, req.Cross)
		if err != nil {
			errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRemove)
			return
		}
	}
	_ = oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_FILE_MANAGER, fmt.Sprintf("用户%v删除文件%v", req.UserName, req.Paths))

	ginutil.Success(ctx, nil)
}

// Remove
//
//	@Summary		文件内容查看接口
//	@Description	文件内容查看接口
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.ReadRequest	true	"入参"
//	@Success		200		{string}	string
//	@Router			/storage/read [post]
func (s *RouteService) Read(ctx *gin.Context) {
	var req dto.ReadRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	fileService, err := s.GetFileService()
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRead)
		return
	}

	if strutil.IsEmpty(req.UserName) {
		req.UserName = ginutil.GetUserName(ctx)
	}

	content, err := fileService.Read(ctx, req.UserName, req.Path, req.Offset, req.Len, req.Cross)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRead)
		return
	}

	ginutil.Success(ctx, content)
}

func (s *RouteService) GetFileService() (service.FileService, error) {
	return s.LocalFileService, nil
}

// Link
//
//	@Summary		下载到我的文件
//	@Description	下载到我的文件
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.LinkRequest	true	"入参"
//	@Success		200
//	@Router			/storage/link [post]
func (s *RouteService) Link(ctx *gin.Context) {
	var req dto.LinkRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	if strutil.IsEmpty(req.UserName) {
		req.UserName = ginutil.GetUserName(ctx)
	}

	fileService, err := s.GetFileService()
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailLink)
		return
	}

	err = fileService.HardLink(ctx, &req)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailLink)
		return
	}

	allSrcPaths := collectionutil.MergeSlice(req.SrcFilePaths, req.SrcDirPaths)
	content := fmt.Sprintf("用户%v链接文件%v => [%v]", req.UserName, allSrcPaths, req.DstPath)
	if req.Cross && !strings.HasPrefix(allSrcPaths[0], req.UserName) {
		pathSlice := strings.Split(allSrcPaths[0], "/")
		filepath.Join("/", ".", req.DstPath)
		content = fmt.Sprintf("用户%v保存%v分享的文件[%v] => [%v]", req.UserName, pathSlice[0], pathSlice[len(pathSlice)-1], strings.Replace(req.DstPath, req.UserName, "home", 1))
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_FILE_MANAGER, content)

	ginutil.Success(ctx, nil)
}

// GenerateAndSendShareCode
//
//	@Summary		生成并发送分享码
//	@Description	生成并发送分享码
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.GenerateShareRequest	true	"入参"
//	@Success		200
//	@Router			/storage/share/send [post]
func (s *RouteService) GenerateAndSendShareCode(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req dto.GenerateShareRequest
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := s.LocalFileService.GenerateShareLink(ctx, req)

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileShareFailed)
		return
	}

	userName := ginutil.GetUserName(ctx)

	sharedUserNames := make([]string, 0)
	reqId := make([]*user.UserIdentity, 0)

	for _, shareUserId := range req.ShareUserList {
		reqId = append(reqId, &user.UserIdentity{Id: shareUserId})
	}
	response, err := client.GetInstance().Users.BatchGetUser(ctx, &user.UserIdentities{
		UserIdentities: reqId,
	})

	for _, userObj := range response.UserObj {
		sharedUserNames = append(sharedUserNames, userObj.Name)
	}

	shareType := "分享"
	if req.ShareType == 1 {
		shareType = "发送"
	}
	content := fmt.Sprintf("用户%v%v文件[%v]给%v", userName, shareType, req.ShareFilePath, sharedUserNames)

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_FILE_MANAGER, content)

	ginutil.Success(ctx, nil)
}

// GetRecordList
//
//	@Summary		获取文件分享记录
//	@Description	获取文件分享记录
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.GetRecordListRequest	true	"入参"
//	@Success		200		{object}	dto.ShareRecordListResponse	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/storage/share/recordList [post]
func (s *RouteService) GetRecordList(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req dto.GetRecordListRequest
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	userID := ginutil.GetUserID(ctx)

	list, err := s.LocalFileService.GetRecordList(ctx, userID, req)

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileShareRecordGetFailed)
	}

	ginutil.Success(ctx, list)
}

// GetShareFile
//
//	@Summary		获取分享文件
//	@Description	获取分享文件
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			code	query		string				true	"分享码"	default("")
//	@Success		200		{object}	dto.ShareFileInfo	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/storage/share/get [get]
func (s *RouteService) GetShareFile(ctx *gin.Context) {
	id := ctx.Query("id")

	if strutil.IsEmpty(id) {
		http.Errf(ctx, errcode.ErrInvalidParam, "id can't empty")
		return
	}

	// 获取分享文件信息
	fileInfo, err := s.LocalFileService.GetShareFile(ctx, snowflake.MustParseString(id))
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileShareGetFailed)
		return
	}

	ginutil.Success(ctx, fileInfo)
}

// UpdateRecordState
//
//	@Summary		修改分享记录状态
//	@Description	修改分享记录状态
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.UpdateRecordRequest	true	"请求参数"
//	@Response		200		"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/storage/share/updateRecordState [put]
func (s *RouteService) UpdateRecordState(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = &dto.UpdateRecordRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	userID := ginutil.GetUserID(ctx)
	if err := s.LocalFileService.UpdateRecordState(ctx, userID, req.RecordIDs); err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileShareRecordUpdateFailed)
		return
	}

	ginutil.Success(ctx, nil)
}

// ReadAll
//
//	@Summary		已读所有分享记录
//	@Description	已读所有分享记录
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Response		200	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/storage/share/readAll [put]
func (s *RouteService) ReadAll(ctx *gin.Context) {
	userID := ginutil.GetUserID(ctx)
	if err := s.LocalFileService.ReadAll(ctx, userID); err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileShareRecordUpdateFailed)
		return
	}

	ginutil.Success(ctx, nil)
}

// ShareCount
//
//	@Summary		获取分享消息数量
//	@Description	获取分享消息数量
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			state	query		string	false	"分享码"	default("")
//	@Success		200		{number}	count	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/storage/share/count [get]
func (s *RouteService) ShareCount(ctx *gin.Context) {
	state := ctx.Query("state")

	var stateI int64
	if strutil.IsNotEmpty(state) {
		stateI, _ = strconv.ParseInt(state, 10, 64)
	}

	count, err := s.LocalFileService.GetShareCount(ctx, ginutil.GetUserID(ctx), int8(stateI))
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileShareRecordGetFailed)
		return
	}

	ginutil.Success(ctx, count)
}
