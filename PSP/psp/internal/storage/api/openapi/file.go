package openapi

import (
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/dto/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/structutil"
)

// PreUpload
//
//	@Summary		文件预上传接口
//	@Description	文件预上传接口
//	@Tags			存储-文件
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			path		formData	string						true	"路径"		default("/")
//	@Param			file_size	formData	int							true	"文件大小"		default(0)
//	@Param			cross		formData	bool						false	"是否跨越用户目录"	default(false)
//	@Param			is_cloud	formData	bool						false	"是否云端"		default(false)
//	@Param			user_name	formData	string						false	"用户名"		default("")
//	@Success		200			{object}	openapi.PreUploadResponse	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/openapi/storage/preUpload [post]
func (s *RouteOpenapiService) PreUpload(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req openapi.PreUploadRequest
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if err := s.validate.Struct(req); err != nil {
		ginutil.Error(ctx, errcode.ErrInvalidParam, err.Error())
		return
	}

	fileService, err := s.GetFileService(req.ComputeType)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRead)
		return
	}

	preUploadResp, err := fileService.PreUpload(ctx, &dto.PreUploadRequest{
		Path:     req.Path,
		FileSize: req.FileSize,
		Cross:    req.IsTemp,
		UserName: ginutil.GetUserName(ctx),
		IsCloud:  openapi.ConvertIsCloud(req.ComputeType),
	})
	if err != nil {
		logger.Errorf("pre upload file err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileUploadInit)
		return
	}

	rsp := &openapi.PreUploadResponse{}
	if err = structutil.CopyStruct(rsp, preUploadResp); err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileUploadInit)
		return
	}

	ginutil.Success(ctx, rsp)
}

// Upload
//
//	@Summary		文件分片上传接口
//	@Description	文件分片上传接口
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	body	openapi.UploadRequest	true	"入参"
//	@Success		200
//	@Router			/openapi/storage/upload [post]
func (s *RouteOpenapiService) Upload(ctx *gin.Context) {

	logger := logging.GetLogger(ctx)

	var req openapi.UploadRequest
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if err := s.validate.Struct(req); err != nil {
		ginutil.Error(ctx, errcode.ErrInvalidParam, err.Error())
		return
	}

	slice := ctx.Request.MultipartForm.File["slice"]

	fileService, err := s.GetFileService(req.ComputeType)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRead)
		return
	}

	userName := ginutil.GetUserName(ctx)
	err = fileService.Upload(ctx, &dto.UploadRequest{
		UploadID:  req.UploadID,
		Path:      req.Path,
		FileSize:  req.FileSize,
		Offset:    req.Offset,
		SliceSize: req.SliceSize,
		Finish:    req.Finish,
		Cross:     req.IsTemp,
		IsCloud:   openapi.ConvertIsCloud(req.ComputeType),
		UserName:  userName,
	}, slice)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileUploadSlice)
		return
	}

	if req.Finish {
		_ = oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_FILE_MANAGER, fmt.Sprintf("【OPENAPI】用户%v上传文件[%v]", userName, req.Path))
	}

	ginutil.Success(ctx, nil)
}

// List
//
//	@Summary		文件列表接口
//	@Description	文件列表接口
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	body	openapi.ListRequest		true	"入参"
//	@Success		200		{array}	openapi.ListResponse	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/openapi/storage/list [post]
func (s *RouteOpenapiService) List(ctx *gin.Context) {
	var req openapi.ListRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	if err := s.validate.Struct(req); err != nil {
		ginutil.Error(ctx, errcode.ErrInvalidParam, err.Error())
		return
	}

	fileService, err := s.GetFileService(req.ComputeType)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileList)
		return
	}

	responses := make([]*openapi.ListResponse, 0)

	// Check the user home path existence
	err = fileService.CheckUserHomePath(ctx, ginutil.GetUserName(ctx), req.IsTemp)
	if err != nil {
		ginutil.Success(ctx, responses)
		return
	}

	err = s.LocalFileService.CheckSharePath(ctx, ginutil.GetUserName(ctx), req.IsTemp)
	if err != nil {
		ginutil.Success(ctx, responses)
		return
	}

	files, err := fileService.List(ctx, ginutil.GetUserName(ctx), req.Path, req.Page, req.IsTemp, req.ShowHideFile, req.FilterRegexpList)
	if err != nil {
		ginutil.Success(ctx, responses)
		return
	}

	files, err = fileService.CheckOnlyReadDir(ctx, req.Path, files)
	if err != nil {
		ginutil.Success(ctx, responses)
		return
	}

	rsp := make([]*openapi.ListResponse, 0)
	if err = structutil.CopyStruct(&rsp, util.ConvertFileRsps(files)); err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileList)
		return
	}

	ginutil.Success(ctx, rsp)
}

func (s *RouteOpenapiService) GetFileService(computeType string) (service.FileService, error) {
	switch computeType {
	case common.Local:
		return s.LocalFileService, nil
	default:
		return nil, status.Error(errcode.ErrInvalidAction, "")
	}
}

// Remove
//
//	@Summary		文件删除接口
//	@Description	文件删除接口
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	body	openapi.RemoveRequest	true	"入参"
//	@Success		200
//	@Router			/openapi/storage/remove [delete]
func (s *RouteOpenapiService) Remove(ctx *gin.Context) {
	var req openapi.RemoveRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	if err := s.validate.Struct(req); err != nil {
		ginutil.Error(ctx, errcode.ErrInvalidParam, err.Error())
		return
	}

	fileService, err := s.GetFileService(req.ComputeType)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRead)
		return
	}

	userName := ginutil.GetUserName(ctx)

	for _, p := range req.Paths {
		path := filepath.Clean(p)

		// Check the path existence
		if !fileService.Exist(ctx, userName, req.IsTemp, path) {
			errcode.ResolveErrCodeMessage(ctx, nil, errcode.ErrFileNotExist)
			return
		}

		// Remove the file
		err := fileService.Remove(ctx, userName, path, req.IsTemp)
		if err != nil {
			errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRemove)
			return
		}
	}
	_ = oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_FILE_MANAGER, fmt.Sprintf("【OPENAPI】用户%v删除文件%v", userName, req.Paths))

	ginutil.Success(ctx, nil)
}

// BatchDownloadPre
//
//	@Summary		批量文件预下载
//	@Description	批量文件预下载接口
//	@Tags			存储-文件
//	@Accept			json
//	@Produce		json
//	@Param			param	body		openapi.BatchDownloadPreRequest	true	"入参"
//	@Success		200		{string}	string
//	@Router			/openapi/storage/batchDownloadPre [post]
func (s *RouteOpenapiService) BatchDownloadPre(ctx *gin.Context) {
	var req openapi.BatchDownloadPreRequest
	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	if err := s.validate.Struct(req); err != nil {
		ginutil.Error(ctx, errcode.ErrInvalidParam, err.Error())
		return
	}

	fileService, err := s.GetFileService(req.ComputeType)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileFailRead)
		return
	}

	token, err := fileService.BatchDownloadPre(ctx, &dto.BatchDownloadPreRequest{
		FilePaths:  req.FilePaths,
		FileName:   req.FileName,
		Cross:      req.IsTemp,
		IsCompress: req.IsCompress,
		IsCloud:    openapi.ConvertIsCloud(req.ComputeType),
		UserName:   ginutil.GetUserName(ctx),
	})
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
//	@Param			param	query	openapi.BatchDownloadRequest	true	"请求参数"
//	@Success		200
//	@Router			/openapi/storage/batchDownload [get]
func (s *RouteOpenapiService) BatchDownload(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &openapi.BatchDownloadRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if err := s.validate.Struct(req); err != nil {
		ginutil.Error(ctx, errcode.ErrInvalidParam, err.Error())
		return
	}

	fileService, err := s.GetFileService(req.ComputeType)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileBatchDownload)
		return
	}

	userName := ginutil.GetUserName(ctx)
	err = fileService.BatchDownload(ctx, userName, req.Token, openapi.ConvertIsCloud(req.ComputeType))
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrFileBatchDownload)
		return
	}
}
